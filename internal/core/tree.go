// internal/core/tree.go

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDirectoryTree(root string, excludeDirs, extensions []string) string {
	var tree []string

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
				ext := filepath.Ext(file.Name())
				for _, validExit := range extensions {
					if ext == validExit {
						tree = append(tree, fmt.Sprintf("%s    |- %s", prefix, file.Name()))
						break
					}
				}
			}
		}
	}

	walk(root, 0)
	return strings.Join(tree, "\n")
}
