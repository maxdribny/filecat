# filecat

A command line tool for combining and analyzing multiple files into a single output file. Written in Go.

## Features

- Combine multiple files into a single output file
- Generate directory tree structure
- Count lines of code
- Copy output to clipboard
- Exclude specific directories
- Support for multiple file extensions

## Installation

```bash
go install github.com/maxdribny/filecat/filecat@latest
```

## Usage

```bash
# Basic usage - combine all .go files
filecat

# Specify multiple file extensions
filecat -ext go,md,txt,cs

# Exclude specific directories
filecat -exclude vendor,tests -ext=go

# Just count lines of code
filecat -count -ext go,md,txt,cs

# Show directory tree
filecat -tree -ext go,cs

# Custom output file
filecat -out output.txt -ext go,cs

# Copy to clipboard
filecat -copy -ext go,cs
```

## Available Options

- `-ext` : Comma seperated list of file extensions to search for (default: "go")
- `-exclude` : Comma-seperated list of directories to exclude
- `-root` : Root directory to start search from (default: current directory)
- `out` : Output file name (default: "combined_files.txt")
- `-count` : Only count lines of code without combining files
- `-tree` : Show directory tree structure
- `-copy` : Copy output to clipboard

## License
MIT