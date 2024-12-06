package main

import "testing"

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
