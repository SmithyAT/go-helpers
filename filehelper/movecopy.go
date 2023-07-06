package filehelper

import (
	"io"
	"os"
	"path/filepath"
)

// MoveFile moves a file from sourcePath to destDir with the filename provided.
// The function takes three parameters:
// sourcePath is the path to the file that needs to be moved
// destDir is the directory to which the file needs to be moved
// filename is the new name of the file once it is moved to the destination directory
//
// MoveFile attempts to open the source file, create the destination directory if it doesn't already exist,
// create a new file in the destination directory, copy the contents from the source file to the dest file,
// and finally delete the source file.
//
// It returns an error if any of the steps fail, otherwise it returns nil indicating successful move.
func MoveFile(sourcePath, destDir, filename string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func(sourcefile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	destPath := filepath.Join(destDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = sourceFile.Close()
	if err != nil {
		return err
	}

	err = os.Remove(sourcePath) // Remove the source file after successful move
	if err != nil {
		return err
	}

	return nil
}

// CopyFile copies a file from sourcePath to a destination directory destDir with a new filename.
// It opens the source file and creates a destination directory if not present.
// Then it creates a new file in the destination directory and copies the source file contents to the new file.
// File handlers are also properly closed after their operations are over.
// If any error occurs during these operations, it will be returned by the function.
func CopyFile(sourcePath, destDir, filename string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func(sourcefile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	destPath := filepath.Join(destDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func(destfile *os.File) {
		_ = destFile.Close()
	}(destFile)

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = sourceFile.Close()
	if err != nil {
		return err
	}

	return nil
}
