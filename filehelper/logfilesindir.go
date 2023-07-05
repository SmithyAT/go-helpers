package filehelper

import (
	"github.com/sirupsen/logrus"
	"io/fs"
	"path/filepath"
)

// LogFilesInDir logs remaining files in a directory
func LogFilesInDir(dirPath string, logPtr *logrus.Entry) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If it is not a directory, then it is a file, print it
		if !d.IsDir() {
			logPtr.Warnf("Unknown file detected %s", path)
		}

		return nil
	})
}
