package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDirectoryTree(root string, excludeDirs, extensions []string) string {
	var tree []string

	// Check if we should match all files
	matchAllFiles := false
	if len(extensions) == 1 && extensions[0] == "" {
		matchAllFiles = true
	}

	var walk func(dir string, level int)
	walk = func(dir string, level int) {
		files, err := os.ReadDir(dir)
		if err != nil {
			return
		}

		prefix := strings.Repeat("    ", level)
		dirName := filepath.Base(dir)

		if level == 0 {
			tree = append(tree, dirName)
		} else {
			tree = append(tree, fmt.Sprintf("%s|- %s", prefix, dirName))
		}

		for _, file := range files {
			path := filepath.Join(dir, file.Name())

			if file.IsDir() {
				skip := false
				for _, excludeDir := range excludeDirs {
					if strings.Contains(path, excludeDir) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				walk(path, level+1)
			} else {
				// Include all files or only files with matching extensions
				if matchAllFiles {
					// Skip hidden files
					if !strings.HasPrefix(file.Name(), ".") {
						tree = append(tree, fmt.Sprintf("%s    |- %s", prefix, file.Name()))
					}
				} else {
					ext := filepath.Ext(file.Name())
					for _, validExt := range extensions {
						if ext == validExt {
							tree = append(tree, fmt.Sprintf("%s    |- %s", prefix, file.Name()))
							break
						}
					}
				}
			}
		}
	}

	walk(root, 0)
	return strings.Join(tree, "\n")
}
