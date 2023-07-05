package filehelper

import (
	"os"
	"path/filepath"
	"strings"
)

// FileInfo represents information about a file.
type FileInfo struct {
	Path string
	Name string
}

// GetFiles scans the provided directory and returns a list of FileInfo structured
// data for each file with the provided file extension. If the file extension is
// empty, the function returns a list of all files in the directory. For any errors
// encountered during filepath walking it returns an error. Here, FileInfo is a
// custom struct that contains information on a file, such as its path and name.
// If an extension is provided that doesn't start with a '.', the function prepends a '.' to it.
// Example:
// GetFiles("/home/user/documents", "txt")
// This would return FileInfo for all text files in /home/user/documents.
func GetFiles(directory, fileExtension string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileExtension == "" {
			if !info.IsDir() {
				files = append(files, FileInfo{Path: path, Name: info.Name()})
			}
		} else {
			// Ensure extensions start with '.'
			if !strings.HasPrefix(fileExtension, ".") {
				fileExtension = "." + fileExtension
			}
			if !info.IsDir() && strings.EqualFold(filepath.Ext(info.Name()), fileExtension) {
				files = append(files, FileInfo{Path: path, Name: info.Name()})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
