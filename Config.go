package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	RetryCount        int    `json:"retryCount"`
	RetryDelay        int    `json:"retryDelay"`
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

// TODO: Make the config global somehow
// TODO: Implement retry count and retry delay
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

	err = ValidateConfig(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func CreateConfig() error {
	config := Config{
		RetryCount:        3,
		RetryDelay:        3,
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

func ValidateConfig(config Config) error {

	if config.RetryCount < 1 {
		return errors.New("retry count must be greater than zero")
	}

	if config.RetryDelay < 1 {
		return errors.New("retry delay must be greater than zero")
	}

	if config.LogDirectory == "" {
		return errors.New("log directory must be set")
	}

	if config.DownloadDirectory == "" {
		return errors.New("download directory must be set")
	}

	if config.ArchiveDirectory == "" {
		return errors.New("archive directory must be set")
	}

	if !IsValidPathName(config.DownloadDirectory) {
		return errors.New("invalid download directory")
	}

	if !IsValidPathName(config.ArchiveDirectory) {
		return errors.New("invalid archive directory")
	}

	if !IsValidPathName(config.LogDirectory) {
		return errors.New("invalid log directory")
	}
	return nil
}
