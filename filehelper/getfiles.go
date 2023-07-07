package filehelper

import (
	"os"
	"path/filepath"
)

// FileInfo represents information about a file.
type FileInfo struct {
	Path string
	Name string
}

// GetFiles scans the provided directory and returns files that match the provided pattern.
// The function returns a slice of FileInfo structures, each representing information about a file,
// and an error when applicable. Each FileInfo contains the file's path and name.
//
// Parameters:
//   - directory: The directory path to scan for files.
//   - pattern: The pattern to match against file names. This uses Go's `filepath.Match`
//     function's pattern syntax.
//
// Returns:
//   - A slice of FileInfo structures. Each FileInfo consists of a Path and a Name.
//   - An error if an error occurs while reading the directory or matching files,
//     otherwise nil.
//
// This function does not look for files in subdirectories of the provided directory,
// only at the topmost level. If the pattern doesn't match any file, this function
// returns an empty slice and nil.
//
// The returned Path in each FileInfo is combining the directory path and the file name.
// The Name in each FileInfo is the name of the file (not including the path).
func GetFiles(directory, pattern string) ([]FileInfo, error) {
	var files []FileInfo

	fileInfo, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	for _, info := range fileInfo {
		if !info.IsDir() {

			if match, _ := filepath.Match(pattern, info.Name()); match {
				files = append(files, FileInfo{Path: filepath.Join(directory, info.Name()), Name: info.Name()})
			}

		}
	}

	return files, nil
}

// GetFilesRecursive is a function that traverses a provided directory recursively
// and returns files matching a given pattern.
//
// It takes two parameters:
// directory - It is the directory where to start the search. This function will recursively search in this directory and all of its subdirectories.
// pattern - The pattern to match file names against. Only the files that match this pattern will be returned.
//
// The function returns a list of FileInfo objects and an error object. FileInfo object includes the full path and the name of the file.
// The error object will be non-nil if there was an error during the function's execution.
//
// Usage:
//
// files, err := GetFilesRecursive("/Users/example_user", "*.go")
//
//	if err != nil {
//	   log.Fatal(err)
//	}
//
//	for _, file := range files {
//	    fmt.Println(file.Path)
//	}
func GetFilesRecursive(directory, pattern string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if match, _ := filepath.Match(pattern, info.Name()); match {
				files = append(files, FileInfo{Path: filepath.Join(directory, info.Name()), Name: info.Name()})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
