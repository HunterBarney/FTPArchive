package main

import (
	"path/filepath"
	"testing"
)

func TestSFTPConnect(t *testing.T) {
	mockClient := &Profile{
		HostName: "test.rebex.net",
		Port:     22,
		Username: "demo",
		Password: "password",
	}

	client, e := connectSFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}
}

func TestSFTPDownloadFile(t *testing.T) {
	mockClient := &Profile{
		HostName:   "test.rebex.net",
		Port:       22,
		Username:   "demo",
		Password:   "password",
		OutputName: "backup_test_sftp",
	}

	client, e := connectSFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}

	remotePath := "pub/example/imap-console-client.png"
	localPath := filepath.Join(mockClient.OutputName, filepath.Base(remotePath))
	e = downloadFileSFTP(client, remotePath, localPath)
	if e != nil {
		t.Fatal("Failed to download: ", e)
	}
}

func TestSFTPDownloadDirectory(t *testing.T) {
	mockClient := &Profile{
		HostName:   "test.rebex.net",
		Port:       22,
		Username:   "demo",
		Password:   "password",
		OutputName: "backup_test_dir_sftp",
	}

	client, e := connectSFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}

	remotePath := "pub"
	localPath := filepath.Join(mockClient.OutputName, filepath.Base(remotePath))
	e = downloadDirectorySFTP(client, remotePath, localPath)
	if e != nil {
		t.Fatal("Failed to download: ", e)
	}
}
