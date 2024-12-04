package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	profilePath := flag.String("profile", "ftptest.json", "The path to the profile.")

	flag.Parse()
	fmt.Println("profilePath:", *profilePath)

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
		e = DisconnectFTP(client)
		if e != nil {
			log.Fatal(e)
		}

	case "SFTP":
		client, e := connectSFTP(&profile)
		if e != nil {
			log.Fatal(e)
		}
		processDownloads(client, &profile)
	default:
		fmt.Println("Unknown protocol")
	}
}
