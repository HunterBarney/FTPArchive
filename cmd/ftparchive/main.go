package main

import (
	"FTPArchive/internal/compression"
	"FTPArchive/internal/config"
	"FTPArchive/internal/ftpclient"
	"FTPArchive/internal/logging"
	"FTPArchive/internal/sftpclient"
	"flag"
	"log"
)

func main() {
	profilePath := flag.String("profile", "/home/hunter/GolandProjects/FTPArchive/test/ftptest.json", "The path to the profile.")

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
		log.Fatal("Error loading profile:", err)
	}

	switch profile.Protocol {
	case "FTP":
		client, e := ftpclient.ConnectFTP(&profile, &configFile)
		if e != nil {
			log.Fatal(e)
		}

		e = ftpclient.ProcessDownloadsFTP(&profile, client, &configFile)
		if e != nil {
			log.Fatal(e)
		}

		e = ftpclient.DisconnectFTP(client)
		if e != nil {
			log.Fatal(e)
		}

	case "SFTP":
		client, e := sftpclient.ConnectSFTP(&profile, &configFile)
		if e != nil {
			log.Fatal(e)
		}
		e = sftpclient.ProcessDownloadsSFTP(client, &profile, &configFile)
		if e != nil {
			log.Fatal(e)
		}
	default:
		log.Println("Unknown protocol")
	}

	e := compression.CompressToZip(profile.OutputName, &configFile)
	if e != nil {
		log.Fatal(e)
	}
}
