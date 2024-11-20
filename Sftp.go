package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
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

func downloadDirectory(client *sftp.Client, remoteDir, localDir string) {
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		fmt.Println(err)
	}

	os.MkdirAll(localDir, 0755)
	for _, file := range files {
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		remoteFilePath = filepath.ToSlash(remoteFilePath)
		localFilePath := filepath.Join(localDir, file.Name())

		fmt.Println("Downloading", remoteFilePath)
		if file.IsDir() {
			downloadDirectory(client, remoteFilePath, localFilePath)
		} else {
			downloadFile(client, remoteFilePath, localFilePath)
		}
	}
}

// downloadFile downloads a single file from the SFTP server
func downloadFile(client *sftp.Client, remotePath, localPath string) {
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		fmt.Println(err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		fmt.Println(err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println(err)
	}
}

func ProcessDownloads(client *sftp.Client, profile *Profile) error {
	for _, item := range profile.Downloads {
		remotePath := item
		localPath := filepath.Join(profile.OutputName, filepath.Base(item))

		stat, err := client.Stat(remotePath)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			downloadDirectory(client, remotePath, localPath)
		} else {
			downloadFile(client, remotePath, localPath)
		}
	}

	return nil
}
