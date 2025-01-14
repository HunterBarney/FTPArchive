package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ConnectFTP takes in a profile and returns an active FTP connection.
func ConnectFTP(profile *Profile, config *Config) (*ftp.ServerConn, error) {
	connectionString := profile.HostName + ":" + strconv.Itoa(profile.Port)
	var err error
	var client *ftp.ServerConn
	// Try to establish a connection to the remote site
	log.Printf("Connecting to: %s", connectionString)
	for i := 1; i < config.RetryCount; i++ {

		client, err = ftp.Dial(connectionString, ftp.DialWithTimeout(5*time.Second))
		// Success!
		if err == nil {
			log.Println("Connected successfully to remote site")
			break
		}

		// Failed, but will retry
		if i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
			continue
		}
	}

	// Connection failed after max retry
	if client == nil {
		return nil, fmt.Errorf("failed to connect to remote site after %d retries", config.RetryCount)
	}

	// Next, log user in.
	log.Printf("Logging in user: %s", profile.Username)
	for i := 1; i < config.RetryCount; i++ {

		err = client.Login(profile.Username, profile.Password)

		// Success!
		if err == nil {
			log.Println("Logged in succesfully to remote site")
			break
		}

		// Failed, but will retry
		if i < config.RetryCount {
			log.Println("Failed. Retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		if i == config.RetryCount {
			return nil, fmt.Errorf("failed to connect to remote site after %d retries", config.RetryCount)
		}
	}

	// Everything was a success!
	return client, nil
}

func DisconnectFTP(client *ftp.ServerConn) error {
	log.Println("Disconnecting from FTP...")
	err := client.Quit()
	if err != nil {
		return err
	}
	log.Println("Successfully disconnected from FTP")
	return nil
}

// DownloadDirectoryFTP recursively downloads all files from the provided remote directory.
func DownloadDirectoryFTP(client *ftp.ServerConn, remoteDir, localDir string, config *Config) error {
	var err error
	var entries []*ftp.Entry

	log.Printf("Getting list of files from directory %s\n", remoteDir)
	for i := 1; i < config.RetryCount; i++ {
		// Get a list of all files in the remote directory
		entries, err = client.List(remoteDir)

		// Success!
		if err == nil {
			break
		}

		// Failed
		if err != nil && i < config.RetryCount {
			log.Println("Failed. retrying...")
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		// Failed, no more retries
		if err != nil && i == config.RetryCount {
			return fmt.Errorf("error getting list of files in remote directory %s after %d retries. err: %s", remoteDir, config.RetryCount, err)
		}
	}

	// Make a local copy of the directory
	log.Printf("Making local copy of directory %s\n", remoteDir)
	err = os.MkdirAll(localDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", localDir, err)
	}

	for _, entry := range entries {
		remotePath := filepath.Join(remoteDir, entry.Name)
		localPath := filepath.Join(localDir, entry.Name)

		if entry.Type == ftp.EntryTypeFolder {
			err = DownloadDirectoryFTP(client, remotePath, localPath, config)
			if err != nil {
				return err
			}
		} else if entry.Type == ftp.EntryTypeFile {
			err = DownloadFileFTP(client, remotePath, localPath, config)
			if err != nil {
				return err
			}
		}
	}
	log.Printf("Directory downloaded successfully: %s\n", localDir)
	return nil
}

// DownloadFileFTP downloads a single file from a remote site
func DownloadFileFTP(client *ftp.ServerConn, remotePath, localPath string, config *Config) error {

	var err error
	var resp *ftp.Response

	log.Printf("Downloading file %s\n", remotePath)
	for i := 1; i < config.RetryCount; i++ {
		// Get file data from FTP site.
		resp, err = client.Retr(remotePath)
		// Success!
		if err == nil {
			break
		}

		if err != nil && i < config.RetryCount {
			log.Println("Failed. retrying...")
			continue
		}

		if err != nil && i == config.RetryCount {
			return fmt.Errorf("unable to read remote file %s. Error %s", remotePath, err)
		}
	}
	defer resp.Close()

	// Creates directory for file if it does not already exist
	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to make local directory for file %s. Error: %w", localPath, err)
	}

	// Create local file
	log.Printf("Creating local file")
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create local file %s. Error: %w", localPath, err)
	}
	defer localFile.Close()

	// Copies remote data to local file
	log.Printf("Copying data...")
	for i := 1; i < config.RetryCount; i++ {
		_, err = io.Copy(localFile, resp)

		// Success!
		if err == nil {
			break
		}

		// Failed. Retries left
		if err != nil && i < config.RetryCount {
			log.Println("Failed. retrying...")
			continue
		}

		// Failed. No retries left.
		if err != nil && i == config.RetryCount {
			return fmt.Errorf("unable to copy data to local file %s. Error: %w", localFile.Name(), err)
		}
	}

	log.Printf("File downloaded successfully: %s\n", localPath)
	return nil
}

// ProcessDownloadsFTP downloads all directories/files from the given profile.
func ProcessDownloadsFTP(profile *Profile, client *ftp.ServerConn, config *Config) error {
	//Create output folder in download directory
	e := os.MkdirAll(profile.OutputName, 0755)
	if e != nil {
		return fmt.Errorf("unable to make output directory for file %s. Error: %w", profile.OutputName, e)
	}

	for _, item := range profile.Downloads {

		remotePath := item
		localPath := filepath.Join(profile.OutputName, filepath.Base(item))
		fileInfo, e := client.GetEntry(item)
		if e != nil {
			return fmt.Errorf("unable to get file info for %s. Error: %w", remotePath, e)
		}
		if fileInfo.Type == ftp.EntryTypeFolder {
			e := DownloadDirectoryFTP(client, remotePath, localPath, config)
			if e != nil {
				return e
			}
		}
		if fileInfo.Type == ftp.EntryTypeFile {
			e := DownloadFileFTP(client, remotePath, localPath, config)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
