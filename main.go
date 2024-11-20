package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	profilePath := flag.String("profile", "profile.json", "The path to the profile.")

	flag.Parse()
	fmt.Println("profilePath:", *profilePath)

	profile, err := LoadProfile(*profilePath)
	if err != nil {
		log.Fatal("Error loading profile:", err)
	}

	switch profile.Protocol {
	case "FTP":
		fmt.Println("FTP Connection")

	case "SFTP":
		client, e := ConnectSFTP(&profile)
		if e != nil {
			log.Fatal(e)
		}
		err := ProcessDownloads(client, &profile)
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Unknown protocol")
	}
}
