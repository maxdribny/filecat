// cmd/filecat/main.go

package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/maxdribny/filecat/internal/core"
	"github.com/spf13/cobra"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "filecat",
		Short: "A tool to combine and analyze source files",
		Long:  `filecat is a command line tool that helps you combine multiple source files into one, generate directory trees, and analyze code files.`,
		RunE:  run,
	}

	rootCmd.Flags().StringP("ext", "e", "go", "Comma-seperated list of file extensions to search for (e.g.)")
	rootCmd.Flags().StringP("exclude", "x", "", "Comma-seperated list of directories to exclude")
	rootCmd.Flags().StringP("root", "r", ".", "Root directory to start search from")
	rootCmd.Flags().StringP("out", "o", "combined_files.txt", "Output file name")
	rootCmd.Flags().BoolP("count", "c", false, "Only count lines of code, don't combine files")
	rootCmd.Flags().Bool("no-combine", false, "Skip combining files into output")
	rootCmd.Flags().BoolP("tree", "t", false, "Show directory tree")
	rootCmd.Flags().BoolP("copy", "y", false, "Copy output to clipboard")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorStyle.Render(err.Error()))
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	config, err := core.NewConfig(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", infoStyle.Render(fmt.Sprintf("Searching for %v", config.FileExtensions)))
	fmt.Printf("%s\n", infoStyle.Render(fmt.Sprintf("Excluding directories: %v", config.ExcludeDirs)))

	// Find all matching files
	files, err := core.FindFiles(config)
	if err != nil {
		return fmt.Errorf("error finding files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found with extensions: %v", config.FileExtensions)
	}

	totalLines := 0
	for _, file := range files {
		totalLines += file.LineCount
	}

	if config.ShowTree {
		tree := core.GenerateDirectoryTree(config.RootDir, config.ExcludeDirs, config.FileExtensions)
		fmt.Println("\nDirectory Structure:")
		fmt.Println("=====================")
		fmt.Println(tree)
	}

	// Always display count if -c/--count flag is set
	if config.CountLines {
		fmt.Println(successStyle.Render(
			fmt.Sprintf("Found %d files with a total of %d lines of code", len(files), totalLines)))
	}

	// Skip combining files if --no-combine is set
	if config.NoCombine {
		return nil
	}

	// Combine files into output
	if err := core.CombineFiles(files, config); err != nil {
		return fmt.Errorf("error combining files: %w", err)
	}

	// Copy to clipboard if option specified
	if config.CopyOutput {
		if err := core.CopyToClipboard(config.OutputFile); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error copying to clipboard: %v", err)))
		} else {
			fmt.Println(successStyle.Render("Content copied to clipboard"))
		}
	}

	fmt.Println(successStyle.Render(
		fmt.Sprintf("Combined %d files into %s with a total of %d lines of code.", len(files), config.OutputFile, totalLines)))

	return nil
}
