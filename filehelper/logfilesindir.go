package filehelper

import (
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

// LogFilesInDir is a function that logs the names of all non-directory files in the provided directory.
// It requires two parameters: a string that represents the path of the directory and a pointer to a logrus.Entry instance for logging.
// It returns an error if it encounters a problem while reading the specified directory.
// For each file in the directory that is not another directory, it will log a warning message that includes the file's name.
// If the function execution is successful, it will return nil, indicating no error occurred.
func LogFilesInDir(dirPath string, logPtr *logrus.Entry) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			logPtr.Warnf("unknown file detected %s", filepath.Join(dirPath, entry.Name()))
		}
	}

	return nil
}

// LogFilesInDirRecursive is a function that logs the presence of non-directory files in a directory recursively.
// It takes the path of the directory to search (dirPath) and a pointer to a logrus.Entry object (logPtr).
//
// This function uses the WalkDir function from the filepath package to traverse the directory tree.
// For each directory entry encountered during the traversal, a function is invoked.
// This function checks if there's any error.
// If there's error, it returns the error immediately.
// Otherwise, it checks if the directory entry is not a directory(i.e., it's a file), and if so,
// it logs a warning with the message "Unknown file detected" along with the file path.
//
// The function returns an error if it failed to traverse the directory or encountered any other error during the process.
//
// Parameters:
//
//	dirPath: a string representing the path of the directory.
//	logPtr: a pointer to a logrus.Entry object used for logging.
//
// Returns:
//   - nil if the directory was successfully traversed and all files within were logged without issue.
//   - error if there was a problem traversing the directory, or if there were any other errors.
func LogFilesInDirRecursive(dirPath string, logPtr *logrus.Entry) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If it is not a directory, then it is a file, print it
		if !d.IsDir() {
			logPtr.Warnf("unknown file detected %s", path)
		}

		return nil
	})
}
