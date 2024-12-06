package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
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

func DownloadDirectoryFTP(entry *ftp.Entry) error {
	fmt.Println("Downloading directory: ", entry.Name)
	return nil
}

func DownloadFileFTP(entry *ftp.Entry) error {
	fmt.Println("Downloading file: ", entry.Name)
	return nil
}

func ProcessDownloadsFTP(profile *Profile, client *ftp.ServerConn) error {
	for _, item := range profile.Downloads {
		file, err := client.GetEntry(item)
		if err != nil {
			return err
		}
		if file.Type == ftp.EntryTypeFolder {
			e := DownloadDirectoryFTP(file)
			return e
		}
		if file.Type == ftp.EntryTypeFile {
			e := DownloadFileFTP(file)
			return e
		}
	}
	return nil
}
