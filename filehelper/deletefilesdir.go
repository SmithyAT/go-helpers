package filehelper

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DeleteFilesInDir is a function that deletes files with the specified extensions
// in the given directory. It takes two parameters:
// 'dirPath', a string indicating the path of directory to be processed,
// and 'extensions', a slice of strings indicating the file extensions to be deleted.
//
// With these parameters, the function will iterate over every item in the directory
// specified by 'dirPath'. If the item is a file (not a directory) and its extension
// matches any in the 'extensions' slice, it will be deleted.
//
// Unlike the DeleteFilesInDirRecursive function, DeleteFilesInDir does not walk through
// subdirectories of the specified directory. This means it only operates on the immediate
// contents of the directory specified by 'dirPath', ignoring any nested directories.
//
// In case of a failure during file deletion, the function will return an error
// detailing the cause of the failure.
//
// It's important to use this function responsibly, considering it will permanently
// delete files.
func DeleteFilesInDir(dirPath string, extensions []string) error {
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		// Ensure extensions start with '.'
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		extMap[strings.ToLower(ext)] = true
	}

	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, d := range dir {
		if d.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if _, ok := extMap[ext]; ok {
			err := os.Remove(filepath.Join(dirPath, d.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteFilesInDirRecursive is a function that deletes files within a specified directory if the file extensions match
// those provided in the 'extensions' argument. It uses 'filepath.WalkDir' to traverse
// the directory tree and 'os.Remove' to delete the files.
//
// dirPath parameter is a string that represents the directory path where the function
// will begin its search.
//
// extensions parameter is a string slice containing the file extensions
// for files that should be deleted. These extensions should start with ".",
// otherwise a "." is prepended.
//
// The function returns an error if any occurs during directory traversal or file deletion.
// If no error occurs, the function returns nil indicating successful file deletion.
func DeleteFilesInDirRecursive(dirPath string, extensions []string) error {
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
