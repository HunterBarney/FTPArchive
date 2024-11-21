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

func connectSFTP(profile *Profile) (*sftp.Client, error) {
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

func downloadDirectory(client *sftp.Client, remoteDir, localDir string) error {
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		fmt.Println(err)
	}

	err = os.MkdirAll(localDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory: ", localDir)
		return err
	}
	for _, file := range files {
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		remoteFilePath = filepath.ToSlash(remoteFilePath)
		localFilePath := filepath.Join(localDir, file.Name())

		fmt.Println("Downloading", remoteFilePath)
		if file.IsDir() {
			err = downloadDirectory(client, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		} else {
			err = downloadFile(client, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadFile(client *sftp.Client, remotePath, localPath string) error {
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		fmt.Println("Error opening remote file: ", remoteFile)
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		fmt.Println("Error creating local file: ", localFile)
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println("Error downloading file: ", remoteFile)
		return err
	}
	return nil
}

func processDownloads(client *sftp.Client, profile *Profile) {
	for _, item := range profile.Downloads {
		remotePath := item
		localPath := filepath.Join(profile.OutputName, filepath.Base(item))

		stat, err := client.Stat(remotePath)
		if err != nil {
			fmt.Println("Error statting remote file: ", remotePath)
			fmt.Println(err)
		}

		if err == nil {
			if stat.IsDir() {
				err = downloadDirectory(client, remotePath, localPath)
				if err != nil {
					fmt.Println("Error downloading directory: ", remotePath)
					fmt.Println(err)
				}
			} else {
				err = downloadFile(client, remotePath, localPath)
				if err != nil {
					fmt.Println("Error downloading file: ", remotePath)
					fmt.Println(err)
				}
			}
		}
	}
}
