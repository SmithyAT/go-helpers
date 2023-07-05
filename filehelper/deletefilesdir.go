package filehelper

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DeleteFilesInDir deletes files in a directory with the given extensions
func DeleteFilesInDir(dirPath string, extensions []string) error {
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		// Ensure extensions start with '.'
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		extMap[strings.ToLower(ext)] = true
	}

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If it is not a directory, then it is a file, delete it if its extension is in the map
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if _, ok := extMap[ext]; ok {
				err := os.Remove(path)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
