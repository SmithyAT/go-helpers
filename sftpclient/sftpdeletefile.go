package sftpclient

import "github.com/pkg/sftp"

// SftpDeleteFile function is used to delete a file at a specified path using an SFTP client.
// It takes an sftp.Client type instance and a path string as parameters.
//
// Returns an error if any issue occurs while deleting the file.
//
//	client: an instance of an SFTP client
//	path: a string representing the path of the file in the SFTP server
//
// Note: The function makes use of the Remove method of the sftp.Client which is responsible for
// removing a file located at the specific path.
//
// Usage:
//
//	err := SftpDeleteFile(client, "/path/to/file")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// If successfully executed, the function will delete the specified file in the SFTP server.
//
//goland:noinspection GoUnusedExportedFunction
func SftpDeleteFile(client *sftp.Client, path string) error {
	return client.Remove(path)
}
