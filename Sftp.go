package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"strconv"
)

func ConnectSFTP(profile *Profile) (*sftp.Client, error) {
	fmt.Println("Connecting to SFTP site: ", profile.HostName)

	sshConfig := &ssh.ClientConfig{
		User: profile.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(profile.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sftpConnection, err := ssh.Dial("tcp", profile.HostName+":"+strconv.Itoa(profile.Port), sshConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sftpConnection)
	if err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func ProcessDownloadsSftp(sftpClient *sftp.Client, profile *Profile) {
	for _, file := range profile.Downloads {
		remoteFile, err := sftpClient.Stat(file)
		if err != nil {
			fmt.Println("Error opening remote file:", err)
		}

		test := remoteFile.IsDir()
		fmt.Println("File: ", file)
		fmt.Println("Is dir: ", test)
	}

}
