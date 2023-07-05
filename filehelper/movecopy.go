package filehelper

import (
	"io"
	"os"
	"path/filepath"
)

// MoveFile moves the file from the source directory to the destination directory.
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

// CopyFile copies the file from the source directory to the destination directory.
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
