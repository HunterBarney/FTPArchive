package main

import (
	"flag"
	"log"
)

func main() {
	profilePath := flag.String("profile", "ftptest.json", "The path to the profile.")

	logFile := initLogging()
	defer logFile.Close()

	flag.Parse()
	log.Println("profilePath:", *profilePath)

	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := LoadProfile(*profilePath)
	if err != nil {
		log.Fatal("Error loading profile:", err)
	}

	switch profile.Protocol {
	case "FTP":
		client, e := ConnectFTP(&profile)
		if e != nil {
			log.Fatal(e)
		}

		e = ProcessDownloadsFTP(&profile, client)
		if e != nil {
			log.Fatal(e)
		}

		e = DisconnectFTP(client)
		if e != nil {
			log.Fatal(e)
		}

	case "SFTP":
		client, e := connectSFTP(&profile)
		if e != nil {
			log.Fatal(e)
		}
		e = processDownloadsSFTP(client, &profile)
		if e != nil {
			log.Fatal(e)
		}
	default:
		log.Println("Unknown protocol")
	}

	e := CompressToZip(profile.OutputName)
	if e != nil {
		log.Fatal(e)
	}
}
