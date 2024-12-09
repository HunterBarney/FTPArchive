package main

import (
	"path/filepath"
	"testing"
)

func TestFTPConnect(t *testing.T) {
	mockClient := &Profile{
		HostName: "test.rebex.net",
		Port:     21,
		Username: "demo",
		Password: "password",
	}

	client, e := ConnectFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}

	e = DisconnectFTP(client)
	if e != nil {
		t.Fatal("Failed to disconnect: ", e)
	}
}

func TestFTPDownloadFile(t *testing.T) {
	mockClient := &Profile{
		HostName:   "test.rebex.net",
		Port:       21,
		Username:   "demo",
		Password:   "password",
		OutputName: "backup_test",
	}

	client, e := ConnectFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}

	remotePath := "pub/example/imap-console-client.png"
	localPath := filepath.Join(mockClient.OutputName, filepath.Base(remotePath))
	e = DownloadFileFTP(client, remotePath, localPath)
	if e != nil {
		t.Fatal("Failed to download: ", e)
	}

	e = DisconnectFTP(client)
	if e != nil {
		t.Fatal("Failed to disconnect: ", e)
	}
}

func TestFTPDownloadDirectory(t *testing.T) {
	mockClient := &Profile{
		HostName:   "test.rebex.net",
		Port:       21,
		Username:   "demo",
		Password:   "password",
		OutputName: "backup_test_dir",
	}

	client, e := ConnectFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}

	remotePath := "pub"
	localPath := filepath.Join(mockClient.OutputName, filepath.Base(remotePath))
	e = DownloadDirectoryFTP(client, remotePath, localPath)
	if e != nil {
		t.Fatal("Failed to download: ", e)
	}

	e = DisconnectFTP(client)
	if e != nil {
		t.Fatal("Failed to disconnect: ", e)
	}
}
