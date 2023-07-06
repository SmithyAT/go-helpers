package filehelper

import (
	"archive/tar"
	"compress/gzip"
	"github.com/sirupsen/logrus"
	"github.com/smithyat/go-helpers/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractAllTarGzInDirectory reads the specified 'srcDir' and for every .tar.gz file found,
// it extracts the file to the 'destDir'. It logs any error encountered during the file operations
// such as opening a file, extracting, and removing a file via the provided logrus Logger.
// It does not look into any subdirectories of 'srcDir'.
//
// The srcDir parameter is the source directory from where the .tar.gz files are to be read.
//
// The destDir parameter is the destination directory where the .tar.gz files are to be extracted.
//
// The logPtr parameter is a pointer to a logrus Logger that logs any file operation errors encountered.
//
// Note: This function does not return any value. All the errors are logged and the function
// moves to the next file operation when an error is encountered.
func ExtractAllTarGzInDirectory(srcDir, destDir string, logPtr *logrus.Entry) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		// Handle error if needed, e.g. srcDir does not exist
		logPtr.Errorf("failed to read directory %s: %v [%s]", srcDir, err, logger.Trace())
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tar.gz") {
			path := filepath.Join(srcDir, file.Name())
			tarGzFile, err := os.Open(path)
			if err != nil {
				logPtr.Errorf("failed to open file %s: %v [%s]", path, err, logger.Trace())
				continue
			}

			logPtr.Infof("extracting file %s", path)
			err = ExtractTarGz(tarGzFile, destDir)
			_ = tarGzFile.Close()

			if err != nil {
				logPtr.Errorf("failed to extract file %s: %v [%s]", path, err, logger.Trace())
				continue
			}

			logPtr.Infof("extracted file %s", path)

			err = os.Remove(path)
			if err != nil {
				logPtr.Errorf("failed to remove file %s: %v [%s]", path, err, logger.Trace())
			}
		}
	}
}

// ExtractAllTarGzInDirectoryRecursive walks through the provided source directory
// recursively and extracts all .tar.gz files found into the destination directory.
// Any errors encountered during the walk or extraction process are logged using
// the provided logrus entry, if it is not nil.
//
// It first opens the .tar.gz file, extracts it using the ExtractTarGz function,
// then deletes the original .tar.gz file. The ExtractTarGz function needs to be
// properly defined in the current or imported package.
//
// This function does not return an error, instead, it logs all errors and returns
// nil to filepath.Walk for graceful failure. Log messages and errors are
// supplemented with trace information from logger.Trace().
//
// Note: The function may fail to extract files or delete the original files if they
// are not accessible due to permission issues or some other constraints.
//
// Parameters:
// srcDir: Source directory path where the function will search for .tar.gz files.
// destDir: Destination directory path where the files will be extracted.
// logPtr: Pointer to a logrus.Entry where log messages will be written.
//
//	It can be nil.
//
// Example:
// ExtractAllTarGzInDirectoryRecursive("/path/to/src", "/path/to/dest", logrus.NewEntry(logger))
func ExtractAllTarGzInDirectoryRecursive(srcDir, destDir string, logPtr *logrus.Entry) {
	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if logPtr != nil {
				logPtr.Errorf("failed to walk path %s: %v [%s]", path, err, logger.Trace())
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".tar.gz") {
			tarGzFile, err := os.Open(path)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("failed to open file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}
			defer func(targzfile *os.File) {
				_ = tarGzFile.Close()
			}(tarGzFile)

			logger.Log.Infof("extracting file %s", path)
			err = ExtractTarGz(tarGzFile, destDir)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("failed to extract file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}
			logger.Log.Infof("extracted file %s", path)

			err = os.Remove(path)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("failed to remove file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}

		}
		return nil
	})
}

// ExtractTarGz takes an io.Reader representing a compressed tarball and a destination string representing
// the directory path where the files should be extracted. This function reads the gzipStream, decompresses it,
// and then extracts each individual file in the tar archive to the location specified by the 'dest' string.
//
// It returns an error if there is any issue during the decompression or extraction process, such as an issue
// creating directories or files, or a problem with the tar archive itself.
//
// The extraction process strictly adheres to the structure of the tarball - directories are created in the 'dest'
// directory as needed, and individual files are written to these directories. The file permissions are also
// maintained during the extraction process.
//
// Note: This function does not handle symbolic links, block devices, or other less common file types in the tar archive.
func ExtractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}

			_ = file.Close()
		}
	}
}

// CreateTarGz is a function that creates a .tar.gz archive from a given source directory.
// It takes a path to the target tar-gz file (`myTarGzFile`), a path to the source directory (`mySourcePath`),
// and a relative path for the file (`relPath`) as arguments.
// If the relative path is not provided, it uses the same directory as the file for it.
// The function walks through every file in the source directory, generates a tar header for each file and writes it to the tar.gz file.
// If the file is not a directory, it writes the content of the file to the archive.
// After successfully writing the file to the archive, it deletes the original file.
// This function returns an error if any occurs during the process.
func CreateTarGz(myTarGzFile, mySourcePath, relPath string) error {
	tarGzFile, err := os.Create(myTarGzFile)
	defer func(targzfile *os.File) {
		_ = tarGzFile.Close()
	}(tarGzFile)

	gw := gzip.NewWriter(tarGzFile)
	defer func(gw *gzip.Writer) {
		_ = gw.Close()
	}(gw)

	tw := tar.NewWriter(gw)
	defer func(tw *tar.Writer) {
		_ = tw.Close()
	}(tw)

	// Walk through every file in the folder
	err = filepath.Walk(mySourcePath, func(file string, fi os.FileInfo, err error) error {
		// Generate tar header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return nil
		}

		var relativeFilePath string
		if relPath == "" {
			relativeFilePath, _ = filepath.Rel(mySourcePath, file) // Set relative path to the same directory as the file
		} else {
			relativeFilePath, _ = filepath.Rel(relPath, file) // Set a specific relative path
		}
		header.Name = relativeFilePath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return nil
		}

		// If it's not a directory, write file content
		if !fi.Mode().IsDir() {
			data, err := os.ReadFile(file)
			if err != nil {
				return nil
			}
			if _, err := tw.Write(data); err != nil {
				return nil
			}
			// Delete the file after archiving
			if err := os.Remove(file); err != nil {
				return nil
			}
		}
		return nil
	})
	return err
}
