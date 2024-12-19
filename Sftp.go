package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// connectSFTP takes in a profile and returns an active SFTP connection.
func connectSFTP(profile *Profile) (*sftp.Client, error) {
	log.Println("Connecting to SFTP site: ", profile.HostName)

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

// downloadDirectorySFTP recursively downloads all files from the provided remote directory.
func downloadDirectorySFTP(client *sftp.Client, remoteDir, localDir string) error {
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		return fmt.Errorf("error getting list of files in remote directory %s. Error %w", remoteDir, err)
	}

	err = os.MkdirAll(localDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", localDir, err)
	}

	// For each item in the directory check if it is a directory or file and run the proper download function.
	for _, file := range files {
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		remoteFilePath = filepath.ToSlash(remoteFilePath) // Converts to unix format as filepath.join on windows defaults to windows format.
		localFilePath := filepath.Join(localDir, file.Name())

		log.Println("Downloading", remoteFilePath)
		if file.IsDir() {
			err = downloadDirectorySFTP(client, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		} else {
			err = downloadFileSFTP(client, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		}
	}
	log.Printf("Directory downloaded successfully: %s\n", localDir)
	return nil
}

// downloadFileSFTP downloads a single file from a remote site
func downloadFileSFTP(client *sftp.Client, remotePath, localPath string) error {
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("unable to read remote file %s. Error %s", remotePath, err)
	}
	defer remoteFile.Close()

	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to make local directory for file %s. Error: %w", localPath, err)
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create local file %s. Error: %w", localPath, err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("unable to copy data to local file %s. Error: %w", localFile.Name(), err)
	}
	return nil
}

// processDownloadsSFTP downloads all directories/files from the given profile.
func processDownloadsSFTP(client *sftp.Client, profile *Profile) error {
	for _, item := range profile.Downloads {
		remotePath := item
		localPath := filepath.Join(profile.OutputName, filepath.Base(item))

		stat, err := client.Stat(remotePath)
		if err != nil {
			return fmt.Errorf("unable to get file info for %s. Error: %w", remotePath, err)
		}

		if stat.IsDir() {
			err = downloadDirectorySFTP(client, remotePath, localPath)
			if err != nil {
				return err
			}
		} else {
			err = downloadFileSFTP(client, remotePath, localPath)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
