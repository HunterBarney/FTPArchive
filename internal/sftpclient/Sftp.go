package sftpclient

import (
	"FTPArchive/internal/config"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ConnectSFTP takes in a profile and returns an active SFTP connection.
func ConnectSFTP(profile *config.Profile, config *config.Config) (*sftp.Client, error) {
	log.Println("Connecting to SFTP site: ", profile.HostName)

	var sftpConnection *ssh.Client
	var err error

	sshConfig := &ssh.ClientConfig{
		User: profile.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(profile.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	for i := 1; i < config.RetryCount; i++ {
		sftpConnection, err = ssh.Dial("tcp", profile.HostName+":"+strconv.Itoa(profile.Port), sshConfig)

		// Success
		if err == nil {
			log.Println("Connected successfully to remote site")
			break
		}

		if err != nil && i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
			continue
		}

		if err != nil && i == config.RetryCount {
			return nil, fmt.Errorf("failed to connect to remote site after %d retries", config.RetryCount)
		}
	}
	sftpClient, err := sftp.NewClient(sftpConnection)
	if err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// downloadDirectorySFTP recursively downloads all files from the provided remote directory.
func downloadDirectorySFTP(client *sftp.Client, remoteDir, localDir string, config *config.Config) error {
	var files []os.FileInfo
	var err error

	log.Printf("Getting list of files from directory %s\n", remoteDir)
	for i := 1; i < config.RetryCount; i++ {
		// Gets list of files from remote directory
		files, err = client.ReadDir(remoteDir)

		// Success!
		if err == nil {
			break
		}

		// Failed, retries left
		if err != nil && i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		// Failed. No retries left.
		if err != nil && i == config.RetryCount {
			return fmt.Errorf("error getting list of files in remote directory %s after %d retries. err: %s", remoteDir, config.RetryCount, err)
		}
	}

	log.Printf("Making local copy of directory %s\n", remoteDir)
	err = os.MkdirAll(localDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", localDir, err)
	}

	// For each item in the directory check if it is a directory or file and run the proper download function.
	for _, file := range files {
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		remoteFilePath = filepath.ToSlash(remoteFilePath) // Converts to unix format as filepath.join on windows defaults to windows format.
		localFilePath := filepath.Join(localDir, file.Name())

		if file.IsDir() {
			err = downloadDirectorySFTP(client, remoteFilePath, localFilePath, config)
			if err != nil {
				return err
			}
		} else {
			err = downloadFileSFTP(client, remoteFilePath, localFilePath, config)
			if err != nil {
				return err
			}
		}
	}
	log.Printf("Directory downloaded successfully: %s\n", localDir)
	return nil
}

// downloadFileSFTP downloads a single file from a remote site
func downloadFileSFTP(client *sftp.Client, remotePath, localPath string, config *config.Config) error {
	var remoteFile *sftp.File
	var err error

	log.Printf("Downloading file %s\n", remotePath)
	for i := 1; i < config.RetryCount; i++ {
		// Gets the remote files data
		remoteFile, err = client.Open(remotePath)

		// Success!
		if err == nil {
			break
		}

		// Failed, retries remaining
		if err != nil && i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		// Failed, no retries left.
		if err != nil && i == config.RetryCount {
			return fmt.Errorf("unable to read remote file %s. Error %s", remotePath, err)
		}
	}
	defer remoteFile.Close()

	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to make local directory for file %s. Error: %w", localPath, err)
	}

	log.Printf("Creating local file")
	// Creates the local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create local file %s. Error: %w", localPath, err)
	}
	defer localFile.Close()

	log.Printf("Copying data...")
	for i := 1; i < config.RetryCount; i++ {
		//Copies the remote data to the local file
		_, err = io.Copy(localFile, remoteFile)

		// Success!
		if err == nil {
			break
		}

		// Failed, retries remaining
		if err != nil && i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		// Failed, no retries left
		if err != nil && i == config.RetryCount {
			return fmt.Errorf("unable to copy data to local file %s. Error: %w", localFile.Name(), err)
		}
	}
	log.Printf("File downloaded successfully: %s\n", localPath)
	return nil
}

// ProcessDownloadsSFTP downloads all directories/files from the given profile.
func ProcessDownloadsSFTP(client *sftp.Client, profile *config.Profile, config *config.Config) error {
	for _, item := range profile.Downloads {
		var err error
		var stat os.FileInfo

		// Make output file.
		err = os.MkdirAll(config.DownloadDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to make output directory for file %s. Error: %w", profile.OutputName, err)
		}

		remotePath := item
		localPath := filepath.Join(profile.OutputName, filepath.Base(item))

		for i := 1; i < config.RetryCount; i++ {
			// Gets info about the remote file/directory
			stat, err = client.Stat(remotePath)
			// Success!
			if err == nil {
				break
			}

			// Failed, retries remaining
			if err != nil && i < config.RetryCount {
				log.Println("Failed. Retrying...")
				time.Sleep(time.Duration(config.RetryDelay) * time.Second)
			}

			// Failed, no retries remaining
			if err != nil && i == config.RetryCount {
				return fmt.Errorf("unable to get file info for %s. Error: %w", remotePath, err)
			}
		}

		if stat != nil && stat.IsDir() {
			err = downloadDirectorySFTP(client, remotePath, localPath, config)
			if err != nil {
				return err
			}
		} else {
			err = downloadFileSFTP(client, remotePath, localPath, config)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
