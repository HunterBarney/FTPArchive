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

// TODO: Enforce the following minimums: Retry count(1), Retry delay(1), LogDirectory/DownloadDirectory/ArchiveDirectory(not blank).
// TODO: Make the config global somehow
// TODO: Decide what defaults for non-required fields should be.
// TODO: Validate config file by making sure each field exists in the file.
// TODO: Implement retry count and retry delay
// TODO: Implement Max allowed errors
// TODO: Implement LogDirectory, DownloadDirectory, and ArchiveDirectory
// Leave email related config stuff alone for now.

// LoadConfig loads the programs config data and returns a struct with the supplied info.
func LoadConfig() (Config, error) {
	var config Config

	// If the config does not exist, create it with the default settings
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		err = CreateConfig()
		if err != nil {
			return config, err
		}
	}

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

func CreateConfig() error {
	config := Config{
		RetryCount:        3,
		RetryDelay:        3,
		MaxAllowedErrors:  0,
		LogDirectory:      "logs",
		DownloadDirectory: "downloads",
		ArchiveDirectory:  "archives",
		SendEmail:         true,
		SendLogOverEmail:  true,
		SMTP: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Host:     "",
			Port:     0,
			Username: "",
			Password: "",
		},
	}

	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
