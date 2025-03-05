# filecat

A command line tool for combining and analyzing multiple files into a single output file. Written in Go.

# FileCat 📁 🔍

A command line tool for combining and analyzing multiple files into a single output file. Written in Go.

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxdribny/filecat)](https://goreportcard.com/report/github.com/maxdribny/filecat)

## Features

- 🔍 **Find & Combine Files**: Easily collect all files of specified types across your project
- 🌲 **Directory Trees**: Generate visual directory structures for better context
- 📊 **Code Statistics**: Count lines of code across your project
- 📋 **Clipboard Integration**: Directly copy combined code to your clipboard
- 🚫 **Smart Exclusions**: Automatically excludes common directories like `.git`, `node_modules`
- 🔄 **Flexible Output**: Control exactly how your code is combined and presented

## Installation

### From Releases

Download the latest release for your platform from the [Releases page](https://github.com/maxdribny/filecat/releases).

### Using Go

```bash
go install github.com/maxdribny/filecat/cmd/filecat@latest
```

### Building from Source

```bash
git clone https://github.com/maxdribny/filecat.git
cd filecat
go build -o filecat ./cmd/filecat
```

## Quick Start

```bash
# Combine all .go files in current directory
filecat -e go

# Combine all .js files and generate a directory tree
filecat -e js -t

# Count lines in all .py files without combining them
filecat -e py -c --no-combine

# Combine files and copy directly to clipboard
filecat -e java -y
```

## Usage

```
filecat - Source File Combiner and Analyzer
=============================================

'filecat' is an easy to use command line tool written in Go that helps you combine multiple file sources into one, 
generate directory trees, and analyze code files.

Usage:
  filecat [flags]

Flags:
  -y, --copy           Copy output file contents to clipboard
  -c, --count          Count lines of code and display total
  -e, --ext string     File extension(s) to search for (comma-separated, without dots)
                       Examples: "go" or "java,js,py"
  -x, --exclude string Directories to exclude (comma-separated)
                       Examples: "node_modules,dist" or "test,vendor"
                       Note: .git, .idea, .vscode, node_modules, build, and dist are excluded by default
  -h, --help           help for filecat
      --no-combine     Skip combining files (useful with -c to only count lines)
  -o, --out string     Output file name
                       Example: "combined_code.txt" (default "combined_files.txt")
  -r, --root string    Root directory to start search from
                       Examples: "." (current directory) or "C:\path\to\project\src" (default ".")
  -t, --tree           Show directory tree of matching files
```

## Common Usage Patterns

1. **Find and combine all Go files in current directory:**
   ```bash
   filecat -e go
   ```

2. **Generate directory tree and count lines (without combining):**
   ```bash
   filecat -e java -t -c --no-combine
   ```

3. **Combine files with specific extension from a directory and save to custom file:**
   ```bash
   filecat -e js -r "./src" -o "javascript_code.txt"
   ```

4. **Work with multiple file extensions:**
   ```bash
   filecat -e "js,ts,jsx" -r "./web" -t
   ```

5. **Combine files and copy result to clipboard:**
   ```bash
   filecat -e py -y
   ```

## Using with AI Assistants

FileCat helps you provide rich context to AI assistants like ChatGPT and Claude:

1. Run FileCat to generate a combined file:
   ```bash
   filecat -e go,yaml -t -o my_project.txt
   ```

2. Upload `my_project.txt` to your AI chat or copy/paste its contents
   ```bash
   filecat -e go,yaml -t -y  # Copies directly to clipboard
   ```

3. Ask questions about your code with complete context

## Why FileCat?

- **Save Time**: No more manually copying multiple files
- **Preserve Context**: Directory structure plus file contents gives AI better understanding
- **Focus**: Only include the files that matter for your question
- **Organize**: Well-formatted output makes it easier for both you and the AI

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra)
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)

---

Made with ❤️ by [Max Dribny](https://github.com/maxdribny)
