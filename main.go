package main

import (
	"flag"
	"log"
)

func main() {
	profilePath := flag.String("profile", "ftptest.json", "The path to the profile.")

	flag.Parse()
	log.Println("profilePath:", *profilePath)

	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logFile := initLogging(&config)
	defer logFile.Close()

	profile, err := LoadProfile(*profilePath, &config)
	if err != nil {
		log.Fatal("Error loading profile:", err)
	}

	switch profile.Protocol {
	case "FTP":
		client, e := ConnectFTP(&profile, &config)
		if e != nil {
			log.Fatal(e)
		}

		e = ProcessDownloadsFTP(&profile, client, &config)
		if e != nil {
			log.Fatal(e)
		}

		e = DisconnectFTP(client)
		if e != nil {
			log.Fatal(e)
		}

	case "SFTP":
		client, e := connectSFTP(&profile, &config)
		if e != nil {
			log.Fatal(e)
		}
		e = processDownloadsSFTP(client, &profile, &config)
		if e != nil {
			log.Fatal(e)
		}
	default:
		log.Println("Unknown protocol")
	}

	e := CompressToZip(profile.OutputName, &config)
	if e != nil {
		log.Fatal(e)
	}
}
