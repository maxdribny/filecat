package core

import (
	"github.com/spf13/cobra"
	"strings"
)

type Config struct {
	ExcludeDirs []string
	Extensions  []string
	RootDir     string
	OutputFile  string
	CountOnly   bool
	ShowTree    bool
	CopyOutput  bool
}

func NewConfig(cmd *cobra.Command) (Config, error) {
	var config Config

	// Get flags from Cobra command
	extensions, _ := cmd.Flags().GetString("ext")
	excludes, _ := cmd.Flags().GetString("exclude")
	config.RootDir, _ = cmd.Flags().GetString("root")
	config.OutputFile, _ = cmd.Flags().GetString("out")
	config.CountOnly, _ = cmd.Flags().GetBool("count")
	config.ShowTree, _ = cmd.Flags().GetBool("tree")
	config.CopyOutput, _ = cmd.Flags().GetBool("copy")

	// Parse extensions
	config.Extensions = strings.Split(extensions, ",")
	for i, ext := range config.Extensions {
		if !strings.HasPrefix(ext, ".") {
			config.Extensions[i] = "." + ext
		}
	}

	if excludes != "" {
		config.ExcludeDirs = strings.Split(excludes, ",")
	}

	defaultExcludes := []string{".git", ".idea", ".vscode", "node_modules", "build", "dist"}
	config.ExcludeDirs = append(config.ExcludeDirs, defaultExcludes...)

	return config, nil
}
