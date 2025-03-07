// internal/core/files.go

package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/schollz/progressbar/v3"
)

type FileInfo struct {
	Path      string
	Ext       string
	LineCount int
}

func FindFiles(config Config) ([]FileInfo, error) {
	var files []FileInfo

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Searching files..."),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetRenderBlankState(true),
	)

	// Check if we should match all files (when "none" was specified)
	matchAllFiles := false
	if len(config.FileExtensions) == 1 && config.FileExtensions[0] == "" {
		matchAllFiles = true
	}

	err := filepath.Walk(config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		bar.Add(1)

		if info.IsDir() {
			for _, excludeDir := range config.ExcludeDirs {
				if strings.Contains(path, excludeDir) {
					return filepath.SkipDir

				}
			}
			return nil
		}

		// If we're matching all files or the extension matches
		if matchAllFiles {
			// Skip hidden files
			fileName := filepath.Base(path)
			if strings.HasPrefix(fileName, ".") {
				return nil
			}

			lineCount, err := countLines(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning, Could not count lines in %s: %v\n", path, err)
				return nil
			}

			files = append(files, FileInfo{
				Path:      path,
				Ext:       filepath.Ext(path),
				LineCount: lineCount,
			})
			return nil
		}

		ext := filepath.Ext(path)
		for _, validExt := range config.FileExtensions {
			if ext == validExt {
				lineCount, err := countLines(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning, Could not count lines in %s: %v\n", path, err)
					continue
				}

				files = append(files, FileInfo{
					Path:      path,
					Ext:       ext,
					LineCount: lineCount,
				})
				break
			}
		}
		return nil
	})

	fmt.Println()

	// Sort files by extension and path
	sort.Slice(files, func(i, j int) bool {
		if files[i].Ext != files[j].Ext {
			return files[i].Ext < files[j].Ext
		}
		return files[i].Path < files[j].Path
	})

	return files, err
}

func CombineFiles(files []FileInfo, config Config) error {
	outFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	// WRite the tree directory
	tree := GenerateDirectoryTree(config.RootDir, config.ExcludeDirs, config.FileExtensions)
	fmt.Fprintln(writer, "Directory Structure:")
	fmt.Fprintln(writer, "===================")
	fmt.Fprintln(writer, tree)
	fmt.Fprintln(writer, "\nThe Source File Contents Are Listed Below, Organized by Extension Under "+
		"the Respective Heading")

	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("Combining files..."),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetRenderBlankState(true),
	)

	currentExt := ""
	for _, file := range files {
		bar.Add(1)

		if currentExt != file.Ext {
			currentExt = file.Ext
			fmt.Fprintf(writer, "\n%s Files:\n", strings.ToUpper(currentExt))
			fmt.Fprintln(writer, strings.Repeat("=", len(currentExt)+7))
			fmt.Fprintln(writer)
		}

		fmt.Fprintf(writer, "// %s\n", file.Path)
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read file %s: %v\n", file.Path, err)
			continue
		}

		fmt.Fprintln(writer, string(content))
		fmt.Fprintln(writer)
	}

	fmt.Println()
	return writer.Flush()
}

func countLines(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func CopyToClipboard(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file for clipboard: %w", err)
	}

	return clipboard.WriteAll(string(content))
}

func GenerateContent(files []FileInfo, config Config) (string, error) {
	var builder strings.Builder

	// Write the tree directory
	tree := GenerateDirectoryTree(config.RootDir, config.ExcludeDirs, config.FileExtensions)
	fmt.Fprintln(&builder, "Directory Structure:")
	fmt.Fprintln(&builder, "===================")
	fmt.Fprintln(&builder, tree)
	fmt.Fprintln(&builder, "\nFile Contents:")
	fmt.Fprintln(&builder, "===============\n")

	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("Processing files..."),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetRenderBlankState(true),
	)

	currentExt := ""
	for _, file := range files {
		bar.Add(1)

		if currentExt != file.Ext {
			currentExt = file.Ext
			fmt.Fprintf(&builder, "\n%s Files:\n", strings.ToUpper(currentExt))
			fmt.Fprintln(&builder, strings.Repeat("=", len(currentExt)+7))
			fmt.Fprintln(&builder)
		}

		fmt.Fprintf(&builder, "\\ %s\n", file.Path)
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read file %s: %v\n", file.Path, err)
			continue
		}

		fmt.Fprintln(&builder, string(content))
		fmt.Fprintln(&builder)
	}

	fmt.Println()
	return builder.String(), nil
}

func CopyContentToClipboard(content string) error {
	return clipboard.WriteAll(content)
}
