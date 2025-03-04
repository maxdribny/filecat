package core

import (
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

	// Parse fileExtensions -> adds a leading "." to the file fileExtensions if
	// missing to ensure all file fileExtensions start with "."
	config.FileExtensions = strings.Split(fileExtensions, ",")
	for i, ext := range config.FileExtensions {
		if !strings.HasPrefix(ext, ".") {
			config.FileExtensions[i] = "." + ext
		}
	}

	if excludes != "" {
		config.ExcludeDirs = strings.Split(excludes, ",")
	}

	defaultExcludes := []string{".git", ".idea", ".vscode", "node_modules", "build", "dist"}
	config.ExcludeDirs = append(config.ExcludeDirs, defaultExcludes...)

	return config, nil
}
