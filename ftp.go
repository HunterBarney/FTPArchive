package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func ConnectFTP(profile *Profile) (*ftp.ServerConn, error) {
	connectionString := profile.HostName + ":" + strconv.Itoa(profile.Port)
	fmt.Println("Connecting to: ", connectionString)
	client, err := ftp.Dial(connectionString, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return client, err
	}

	fmt.Println("Logging in user: ", profile.Username)
	err = client.Login(profile.Username, profile.Password)
	if err != nil {
		return client, err
	}

	fmt.Println("Successfully logged in user: ", profile.Username)
	return client, nil
}

func DisconnectFTP(client *ftp.ServerConn) error {
	fmt.Println("Disconnecting from FTP...")
	err := client.Quit()
	if err != nil {
		return err
	}
	fmt.Println("Successfully disconnected from FTP")
	return nil
}

func DownloadDirectoryFTP(client *ftp.ServerConn, remoteDir, localDir string) error {
	entries, err := client.List(remoteDir)
	if err != nil {
		return fmt.Errorf("error getting list of files in remote directory %s. Error %w", remoteDir, err)
	}

	err = os.MkdirAll(localDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", localDir, err)
	}

	for _, entry := range entries {
		remotePath := filepath.Join(remoteDir, entry.Name)
		localPath := filepath.Join(localDir, entry.Name)

		if entry.Type == ftp.EntryTypeFolder {
			err = DownloadDirectoryFTP(client, remotePath, localPath)
			if err != nil {
				return err
			}
		} else if entry.Type == ftp.EntryTypeFile {
			err = DownloadFileFTP(client, remotePath, localPath)
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("Directory downloaded successfully: %s\n", localDir)
	return nil
}

func DownloadFileFTP(client *ftp.ServerConn, remotePath, localPath string) error {
	resp, err := client.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("unable to read remote file %s. Error %s", remotePath, err)
	}
	defer resp.Close()

	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to make local directory for file %s. Error: %w", localPath, err)
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create local file %s. Error: %w", localPath, err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, resp)
	if err != nil {
		return fmt.Errorf("unable to copy data to local file %s. Error: %w", localFile.Name(), err)
	}

	fmt.Printf("File downloaded successfully: %s\n", localPath)
	return nil
}

func ProcessDownloadsFTP(profile *Profile, client *ftp.ServerConn) error {
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
			e := DownloadDirectoryFTP(client, remotePath, localPath)
			if e != nil {
				return e
			}
		}
		if fileInfo.Type == ftp.EntryTypeFile {
			e := DownloadFileFTP(client, remotePath, localPath)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
