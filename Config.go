package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	RetryCount        int    `json:"retryCount"`
	RetryDelay        int    `json:"retryDelay"`
	MaxAllowedErrors  int    `json:"maxAllowedErrors"`
	LogDirectory      string `json:"logDirectory"`
	DownloadDirectory string `json:"downloadDirectory"`
	ArchiveDirectory  string `json:"archiveDirectory"`
	SendEmail         bool   `json:"sendEmail"`
	SMTP              struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	SendLogOverEmail bool `json:"sendLogOverEmail"`
}

//TODO: Make config if it does not exist
//TODO: Decide what defaults for non-required fields should be.
//TODO: Validate config
//TODO: Implement retry count and retry delay
//TODO: Implement Max allowed errors
//TODO: Implement LogDirectory, DownloadDirectory, and ArchiveDirectory
//Leave email related config stuff alone for now.

// LoadConfig loads the programs config data and returns a struct with the supplied info.
func LoadConfig() (Config, error) {
	var config Config

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}
