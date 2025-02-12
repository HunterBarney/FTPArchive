package main

import (
	"FTPArchive/internal/awsclient"
	"FTPArchive/internal/compression"
	"FTPArchive/internal/config"
	"FTPArchive/internal/emailclient"
	"FTPArchive/internal/ftpclient"
	"FTPArchive/internal/gcp"
	"FTPArchive/internal/logging"
	"FTPArchive/internal/sftpclient"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	profilePath := flag.String("profile", "profile.json", "The path to the profile.")

	flag.Parse()
	log.Println("profilePath:", *profilePath)

	configFile, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logFile := logging.InitLogging(&configFile)
	defer logFile.Close()

	profile, err := config.LoadProfile(*profilePath, &configFile)
	if err != nil {
		failState(&profile, &configFile, logFile, err)
	}

	switch profile.Protocol {
	case "FTP":
		client, e := ftpclient.ConnectFTP(&profile, &configFile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}

		e = ftpclient.ProcessDownloadsFTP(&profile, client, &configFile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}

		e = ftpclient.DisconnectFTP(client)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}

	case "SFTP":
		client, e := sftpclient.ConnectSFTP(&profile, &configFile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}
		e = sftpclient.ProcessDownloadsSFTP(client, &profile, &configFile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}
	default:
		failState(&profile, &configFile, logFile, fmt.Errorf("unknown protocol: %s", profile.Protocol))
	}

	e := compression.CompressToZip(profile.OutputName, &configFile)
	if e != nil {
		failState(&profile, &configFile, logFile, e)
	}

	// Handle uploading
	if profile.UploadPlatform == "aws" || profile.UploadPlatform == "AWS" {
		e = awsclient.UploadFileAWS(&profile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}
	} else if profile.UploadPlatform == "gcp" || profile.UploadPlatform == "GCP" {
		e = gcp.UploadArchiveGcp(&profile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}
	} else {
		log.Printf("Unknown upload platform %s", profile.UploadPlatform)
	}

	// Success!! If email is enabled, send the success email
	if configFile.SendEmail {
		body := "Profile of " + profile.OutputName + " has been archived successfully!"
		err = emailclient.SendEmail("FTPArchive Success!", body, &configFile, &profile, logFile)
		if err != nil {
			log.Fatalf("Error sending email: %v", err)
		}
	}

}

func failState(profile *config.Profile, config *config.Config, logFile *os.File, err error) {
	log.Printf(err.Error())
	err = emailclient.SendFailEmail(config, profile, err, logFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}
