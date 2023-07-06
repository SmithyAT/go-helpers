package filehelper

import (
	"compress/gzip"
	"io"
	"os"
)

// GzipFile compresses a file to a gzip file
func GzipFile(FilePath string, outputFilePath string) error {
	// Open the file for reading
	tarFile, err := os.Open(FilePath)
	if err != nil {
		return err
	}
	defer func(tarFile *os.File) {
		_ = tarFile.Close()
	}(tarFile)

	// Create the output file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(outputFile)
	defer func(gzipWriter *gzip.Writer) {
		_ = gzipWriter.Close()
	}(gzipWriter)

	// Copy the tar file content to the gzip writer
	_, err = io.Copy(gzipWriter, tarFile)
	if err != nil {
		return err
	}

	// Make sure to check the error on Close.
	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// UnzipFile decompresses a gzip file
func UnzipFile(inputFilePath string, outputFilePath string) error {
	// Open the gzip file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create a new gzip reader
	gzipReader, err := gzip.NewReader(inputFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy the gzip content to the output file
	_, err = io.Copy(outputFile, gzipReader)
	if err != nil {
		return err
	}

	return nil
}
