package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type Config struct {
	ExcludeDirs    []string
	FileExtensions []string
	RootDir        string
	OutputFile     string
	CountLines     bool
	NoCombine      bool
	ShowTree       bool
	CopyOutput     bool
}

func NewConfig(cmd *cobra.Command) (Config, error) {
	var config Config

	// Get flags from Cobra command
	fileExtensions, _ := cmd.Flags().GetString("ext")
	excludes, _ := cmd.Flags().GetString("exclude")
	config.RootDir, _ = cmd.Flags().GetString("root")
	config.OutputFile, _ = cmd.Flags().GetString("out")
	config.CountLines, _ = cmd.Flags().GetBool("count")
	config.NoCombine, _ = cmd.Flags().GetBool("no-combine")
	config.ShowTree, _ = cmd.Flags().GetBool("tree")
	config.CopyOutput, _ = cmd.Flags().GetBool("copy")

	// If none is specified, match all the files
	if fileExtensions == "none" {
		// Use an empty string as a marker for "match all files"
		config.FileExtensions = []string{""}
	} else if fileExtensions == "" {
		// If no extension is specified, require the user to provide one or use "none"
		return config, fmt.Errorf("no file extensions specified. Use -e/--ext flag to specify extensions or use -e none to match all files")
	} else {
		// Normal case: parse specified extensions
		config.FileExtensions = strings.Split(fileExtensions, ",")
		for i, ext := range config.FileExtensions {
			// Skip empty extensions (could happen with "ext1,,ext2")
			if ext == "" {
				continue
			}

			if !strings.HasPrefix(ext, ".") {
				config.FileExtensions[i] = "." + ext
			}
		}

		// Filter out any empty extensions that might have resulted from the split
		var filteredExtensions []string
		for _, ext := range config.FileExtensions {
			if ext != "" {
				filteredExtensions = append(filteredExtensions, ext)
			}
		}
		config.FileExtensions = filteredExtensions

		if len(config.FileExtensions) == 0 {
			return config, fmt.Errorf("no valid file extensions specified. Use -e/--ext flag to specify extensions or use -e none to match all files")
		}
	}

	if excludes != "" {
		config.ExcludeDirs = strings.Split(excludes, ",")
	}

	defaultExcludes := []string{".git", ".idea", ".vscode", "node_modules", "build", "dist"}
	config.ExcludeDirs = append(config.ExcludeDirs, defaultExcludes...)

	return config, nil
}
