package filehelper

import "os"

// Touch behaves like the 'touch' command in Unix - it creates a file if it doesn't exist,
// and if it does exist, it deletes and recreates the file to update the timestamp.
func Touch(path string) error {
	// Check if file exists
	_, err := os.Stat(path)
	if err == nil {
		// File exists, delete it
		err = os.Remove(path)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		// Some other error occurred
		return err
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return file.Close()
}
