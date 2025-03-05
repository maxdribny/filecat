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

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	//TODO: Make a style for the description of the tool (pref green)
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "filecat",
		Short: "A tool to combine and analyze source files",
		Long: `'filecat'' is an easy to use command line tool written in Go that helps you combine multiple file sources into one, 
generate directory trees, and analyze code files.

Examples:
 # Combine all .go files in the current directory into combined_files.txt
 filecat -e go

 # Combine all .java files from a specific directory, with tree view, into a custom output file
 filecat -e java -r "C:\path\to\project\src" -t -o "combined_java.txt"

 # Only count lines of code for .js files, without combining
 filecat -e js -r "./web/scripts" -c --no-combine

 # Combine all .py files, show directory tree, and copy to clipboard
 filecat -e py -t -y

 # Exclude specific directories when searching for .cpp files
 filecat -e cpp -x "tests,vendor,third_party"`,
		Example: `  filecat -e go
 filecat -e java -r "C:\path\to\project\src" -t -o output.txt
 filecat -e js,ts -r "./web" -c -t`,
		RunE: run,
	}

	// Define flags with improved descriptions
	rootCmd.Flags().StringP("ext", "e", "",
		`File extension(s) to search for (comma-separated, without dots)
Examples: "go" or "java,js,py"`)

	rootCmd.Flags().StringP("exclude", "x", "",
		`Directories to exclude (comma-separated)
Examples: "node_modules,dist" or "test,vendor"
Note: .git, .idea, .vscode, node_modules, build, and dist are excluded by default`)

	rootCmd.Flags().StringP("root", "r", ".",
		`Root directory to start search from
Examples: "." (current directory) or "C:\path\to\project\src"`)

	rootCmd.Flags().StringP("out", "o", "combined_files.txt",
		`Output file name
Example: "combined_code.txt"`)

	rootCmd.Flags().BoolP("count", "c", false,
		`Count lines of code and display total`)

	rootCmd.Flags().Bool("no-combine", false,
		`Skip combining files (useful with -c to only count lines)`)

	rootCmd.Flags().BoolP("tree", "t", false,
		`Show directory tree of matching files`)

	rootCmd.Flags().BoolP("copy", "y", false,
		`Copy output file contents to clipboard`)

	// Custom help template
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(helpStyle.Render("\nfilecat - Source File Combiner and Analyzer"))
		fmt.Println(helpStyle.Render("=============================================\n"))

		// Display the long description
		fmt.Println(cmd.Long)
		fmt.Println()

		cmd.Usage()
		fmt.Println()
		fmt.Println(helpStyle.Render("Common Usage Patterns:"))
		fmt.Println(helpStyle.Render("---------------------"))
		fmt.Println("1. Find and combine all .go files in current directory:")
		fmt.Println("   filecat -e go")
		fmt.Println()
		fmt.Println("2. Generate directory tree and count lines (without combining):")
		fmt.Println("   filecat -e java -t -c --no-combine")
		fmt.Println()
		fmt.Println("3. Combine files with specific extension from a directory and save to custom file:")
		fmt.Println("   filecat -e js -r \"./src\" -o \"javascript_code.txt\"")
		fmt.Println()
		fmt.Println("4. Work with multiple file extensions:")
		fmt.Println("   filecat -e \"js,ts,jsx\" -r \"./web\" -t")
		fmt.Println()
		fmt.Println("5. Combine files and copy result to clipboard:")
		fmt.Println("   filecat -e py -y")
		fmt.Println()
		fmt.Println(helpStyle.Render("Note: Flags can be specified in any order"))
	})

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

	// Handle copy to clipboard with or without combining
	if config.CopyOutput {
		if config.NoCombine {
			// Generate content in memory without creating a file
			content, err := core.GenerateContent(files, config)
			if err != nil {
				return fmt.Errorf("error generating content: %w", err)
			}

			if err := core.CopyContentToClipboard(content); err != nil {
				fmt.Println(errorStyle.Render(fmt.Sprintf("Error copying to clipboard: %v", err)))
			} else {
				fmt.Println(successStyle.Render("Content copied to clipboard"))
			}
		}
	}

	// Skip combining files if --no-combine is set
	if config.NoCombine {
		return nil
	}

	// Combine files into output
	if err := core.CombineFiles(files, config); err != nil {
		return fmt.Errorf("error combining files: %w", err)
	}

	// Copy to clipboard if option specified (and not already done)
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
