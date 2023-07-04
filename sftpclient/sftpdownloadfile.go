package sftpclient

import (
	"github.com/pkg/sftp"
	"io"
	"os"
)

// SftpDownloadFile takes an established sftp client,
// a remote file path on the sftp server, and a local file path on your system.
// The function opens the remote file and creates a new file on the local path.
// It then copies the contents of the remote file into the local file.
// If it encounters errors during any of these operations (opening source file,
// creating destination file, or copying contents), it returns the error,
// including errors occurred during closing either of the files.
//
// client: An established sftp client.
// remotePath: The full path (including filename) of the remote file on the sftp server.
// localPath: The full path on the local system where the file will be created.
//
// Returns an error if any file operations fail, otherwise nil.
//
//goland:noinspection GoUnusedExportedFunction
func SftpDownloadFile(client *sftp.Client, remotePath, localPath string) (retErr error) {
	// Open remote file
	srcFile, err := client.Open(remotePath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := srcFile.Close()
		if closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	// Create local file
	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := dstFile.Close()
		if closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	// Copy contents to local file
	_, err = io.Copy(dstFile, srcFile)
	return err
}
