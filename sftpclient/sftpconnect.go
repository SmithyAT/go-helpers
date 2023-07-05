package sftpclient

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// SftpConnect establishes an SFTP connection to a server.
// It takes in the username and address of the server, as well as the user's password and private key as parameters.
// If the private key is valid, it uses that for authentication. If not, it falls back on password authentication.
// If none of the authentication methods succeed, it returns an error.
// It returns the connected sftp.Client and any errors encountered.
//
// Params:
//
//	user: Username for the server
//	password: User's password
//	addr: Address of the server
//	key: User's private key
//
// Returns:
//
//	*sftp.Client: The connected SFTP client
//	error: Any errors encountered while connecting to the server
func SftpConnect(user, password, addr string, key []byte) (*sftp.Client, error) {
	var authMethod ssh.AuthMethod

	signer, err := ssh.ParsePrivateKey(key)
	if err == nil {
		authMethod = ssh.PublicKeys(signer)
	} else if password != "" {
		authMethod = ssh.Password(password)
	} else {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	// Connect to SSH server with a timeout
	var client *ssh.Client
	for retries := 0; retries < 3; retries++ {
		client, err = ssh.Dial("tcp", addr, config)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				time.Sleep(2 * time.Second) // Wait before retrying
				continue
			} else {
				return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
			}
		} else {
			break
		}
	}

	// If still not able to connect
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}

	// Create new SFTP Client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		closeErr := client.Close()
		if closeErr != nil {
			fmt.Printf("error on closing ssh connection: %v", closeErr)
		}
		return nil, fmt.Errorf("unable to create SFTP client: %v", err)
	}

	return sftpClient, nil
}
