package main

import "testing"

func TestFTPConnect(t *testing.T) {
	mockClient := &Profile{
		HostName: "ftp.dlptest.com",
		Port:     21,
		Username: "dlpuser",
		Password: "rNrKYTX9g7z3RgJRmxWuGHbeu",
	}

	client, e := ConnectFTP(mockClient)
	if e != nil {
		t.Fatal("Failed to connect: ", e)
	}

	if client == nil {
		t.Fatal("Failed to create client")
	}
}
