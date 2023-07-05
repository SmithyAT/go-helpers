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

// ExtractAllTarGzInDirectory extracts all tar.gz files in the given directory to the given destination directory.
func ExtractAllTarGzInDirectory(srcDir, destDir string, logPtr *logrus.Entry) {
	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if logPtr != nil {
				logPtr.Errorf("Failed to walk path %s: %v [%s]", path, err, logger.Trace())
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".tar.gz") {
			tarGzFile, err := os.Open(path)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("Failed to open file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}
			defer func(targzfile *os.File) {
				_ = tarGzFile.Close()
			}(tarGzFile)

			logger.Log.Infof("Extracting file %s", path)
			err = ExtractTarGz(tarGzFile, destDir)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("Failed to extract file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}
			logger.Log.Infof("Extracted file %s", path)

			err = os.Remove(path)
			if err != nil {
				if logPtr != nil {
					logPtr.Errorf("Failed to remove file %s: %v [%s]", path, err, logger.Trace())
				}
				return nil
			}

		}
		return nil
	})
}

// ExtractTarGz extracts a tar.gz file to the given destination directory.
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

// CreateTarGz creates a tar.gz file
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
