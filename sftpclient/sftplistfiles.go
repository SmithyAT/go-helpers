package sftpclient

import (
	"github.com/pkg/sftp"
	"os"
)

// SftpListFiles lists all files at the given path on the SFTP server.
//
// The function accepts two parameters: a `client` which is an
// instance of the sftp.Client and a `path` which is a string representing
// the path in the SFTP server where the files are to be listed.
//
// If the operation is successful, it will return a list of os.FileInfo
// objects, each representing a file or directory on the given path, and
// a nil error. In case of any error while reading the directory, it will
// return a non-nil error along with a nil list.
//
// Usage:
//
//	fileInfoList, err := SftpListFiles(client, "/path/to/directory")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, fileInfo := range fileInfoList {
//	    fmt.Println(fileInfo.Name())
//	}
//
// Parameters:
// client : *sftp.Client -  a pointer to the SFTP client object
// path : string - the directory path
//
// Returns:
// []os.FileInfo - a slice containing os.FileInfo objects
// error - an error object that reports why the operation failed, or nil if it succeeded
func SftpListFiles(client *sftp.Client, path string) ([]os.FileInfo, error) {
	return client.ReadDir(path)
}
