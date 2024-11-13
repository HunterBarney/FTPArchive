package main

import (
	"encoding/json"
	"log"
	"os"
)

type Profile struct {
	HostName string `json:"hostName"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	// Can be SFTP or FTP
	Protocol string `json:"protocol"`
	// List of files/directories to download. * can be used to specify everything.
	Downloads  []string `json:"downloads"`
	OutputName string   `json:"outputName"`
}

func LoadProfile(fileName string) (Profile, error) {
	var profile Profile

	profileFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer profileFile.Close()

	decoder := json.NewDecoder(profileFile)

	err = decoder.Decode(&profile)
	if err != nil {
		log.Fatal(err)
	}

	return profile, nil
}
