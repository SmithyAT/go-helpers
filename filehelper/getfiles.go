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

// GetFiles returns a slice of FileInfo items of a given directory,
// where each FileInfo contains the path and name of each file that meets the specified criteria.
// If a file extension is specified in the fileExtension parameter, then only files with that
// extension will be included in the result.
// If the fileExtension parameter is an empty string, the function will return info about all files, regardless of extension.
// The function only looks directly in the specified directory and does not traverse any subdirectories.
//
// Parameters:
// directory: The directory to scan for files.
// fileExtension: The file extension to filter by. Leave this an empty string
// to include all files irrespective of their extension.
//
// Return:
// Slice of FileInfo: It contains information about all matching files.
// error: The function will return an error if there was an issue reading the directory.
func GetFiles(directory, fileExtension string) ([]FileInfo, error) {
	var files []FileInfo

	fileInfo, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	// Ensure extensions start with '.'
	if !strings.HasPrefix(fileExtension, ".") && fileExtension != "" {
		fileExtension = "." + fileExtension
	}

	for _, info := range fileInfo {
		if !info.IsDir() {
			if fileExtension == "" {
				files = append(files, FileInfo{Path: filepath.Join(directory, info.Name()), Name: info.Name()})
			} else {
				if strings.EqualFold(filepath.Ext(info.Name()), fileExtension) {
					files = append(files, FileInfo{Path: filepath.Join(directory, info.Name()), Name: info.Name()})
				}
			}
		}
	}

	return files, nil
}

// GetFilesRecursive scans the provided directory and returns a list of FileInfo structured
// data for each file with the provided file extension. If the file extension is
// empty, the function returns a list of all files in the directory. For any errors
// encountered during filepath walking it returns an error. Here, FileInfo is a
// custom struct that contains information on a file, such as its path and name.
// If an extension is provided that doesn't start with a '.', the function prepends a '.' to it.
// Example:
// GetFiles("/home/user/documents", "txt")
// This would return FileInfo for all text files in /home/user/documents.
func GetFilesRecursive(directory, fileExtension string) ([]FileInfo, error) {
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
