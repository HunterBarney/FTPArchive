package main

import (
	"FTPArchive/internal/awsclient"
	"FTPArchive/internal/cleanup"
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
	ProgramVersion := "0.4.0 - Beta"

	profilePath := flag.String("profile", "profile.json", "The path to the profile.")
	manualDownload := flag.Bool("download", false, "Runs only the download function, unless another function is explicitly passed as argument")
	manualArchive := flag.Bool("archive", false, "Runs only the archive function unless another function is explicitly passed as argument")
	manualUpload := flag.Bool("upload", false, "Runs only the upload function unless another function is explicitly passed as argument")
	help := flag.Bool("help", false, "Shows this help")
	manualMode := false

	flag.Parse()

	// If any manual flags are set to true above, enable manual mode
	if *manualDownload || *manualArchive || *manualUpload {
		manualMode = true
	}

	if *help {
		fmt.Println("FTPArchive v" + ProgramVersion)
		fmt.Println("Usage: ftparchive [OPTIONS]")
		flag.PrintDefaults()

		os.Exit(0)
	}

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

	if !manualMode || *manualDownload {
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
	}

	if !manualMode || *manualArchive {
		e := compression.CompressToZip(profile.OutputName, &configFile)
		if e != nil {
			failState(&profile, &configFile, logFile, e)
		}
	}

	if !manualMode || *manualUpload {
		// Handle uploading
		if profile.UploadPlatform == "aws" || profile.UploadPlatform == "AWS" {
			e := awsclient.UploadFileAWS(&profile)
			if e != nil {
				failState(&profile, &configFile, logFile, e)
			}
		} else if profile.UploadPlatform == "gcp" || profile.UploadPlatform == "GCP" {
			e := gcp.UploadArchiveGcp(&profile)
			if e != nil {
				failState(&profile, &configFile, logFile, e)
			}
		} else {
			log.Printf("Unknown upload platform %s", profile.UploadPlatform)
		}
	}

	err = cleanup.Cleanup(&profile)
	if err != nil {
		failState(&profile, &configFile, logFile, err)
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
	if profile.CleanupOnFail {
		err = cleanup.Cleanup(profile)
		if err != nil {
			log.Printf("Error running cleanup: %v", err)
		}
	}
	if config.SendEmail {
		err = emailclient.SendFailEmail(config, profile, err, logFile)
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}
