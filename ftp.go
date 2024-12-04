package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"strconv"
	"time"
)

func ConnectFTP(profile *Profile) (*ftp.ServerConn, error) {
	connectionString := profile.HostName + ":" + strconv.Itoa(profile.Port)
	fmt.Print("Connecting to: ", connectionString)
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

func DownloadDirectoryFTP() {

}

func DownloadFileFTP() {

}

func ProcessDownloadsFTP() {

}
