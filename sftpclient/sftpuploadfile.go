package sftpclient

import (
	"github.com/pkg/sftp"
	"io"
	"os"
)

// SftpUploadFile uploads a file to a remote server over SFTP.
// It takes a sftp.Client object, a local file path and the remote file path as arguments.
// The local file at 'localPath' is opened for reading, while a new file at 'remotePath' is created on the remote server.
// The contents of the local file are then copied over to the remote file.
//
// The function handles cleanup of resources, ensuring that the file handles are properly closed once the operation is complete
// (even in cases where an error might have occurred).
//
// If an error occurs at any point during this operation, it is immediately returned to the calling function.
// This includes errors that occur when opening the files, creating the remote file, copying the contents,
// and errors that occur while closing the local or remote file.
//
// The function returns the error value from the last operation (if any), or nil if the operation was successful.
func SftpUploadFile(client *sftp.Client, localPath, remotePath string) (retErr error) {
	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := localFile.Close()
		if closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	// Create remote file
	remoteFile, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := remoteFile.Close()
		if closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	// Copy contents to remote file
	_, err = io.Copy(remoteFile, localFile)
	return err
}
