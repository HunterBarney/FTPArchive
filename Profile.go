package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// LoadProfile reads data from a profile json file and returns a profile object.
func LoadProfile(fileName string, config *Config) (Profile, error) {
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

	newFileName := FormatDateTime(profile.OutputName, time.Now())
	if IsValidPathName(newFileName) {
		profile.OutputName = filepath.Join(config.DownloadDirectory, newFileName)
	} else {
		return profile, errors.New("output name invalid")
	}

	// Can not create an output folder that already exists
	_, err = os.Stat(profile.OutputName)
	if err == nil {
		// File or directory exists
		return profile, errors.New("output file already exists")
	} else if !os.IsNotExist(err) {
		// An unexpected error occurred
		return profile, err
	}

	return profile, nil
}

// IsValidPathName checks the provided string for invalid characters
func IsValidPathName(path string) bool {
	invalidChars := []string{"/", "<", ">", "\"", "\\", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return false
		}
	}

	return true
}

func IsValidNameConfig(configPath string) bool {
	invalidChars := []string{"<", ">", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(configPath, char) {
			return false
		}
	}

	return true
}
