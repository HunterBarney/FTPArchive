package main

import (
	"flag"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
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
		ConnectSFTP(&profile)

	default:
		fmt.Println("Unknown protocol")
	}
}

func ConnectSFTP(profile *Profile) {
	fmt.Println("Connecting to SFTP site: ", profile.HostName)

	sshConfig := &ssh.ClientConfig{
		User: profile.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(profile.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sftpConnection, err := ssh.Dial("tcp", profile.HostName+":"+strconv.Itoa(profile.Port), sshConfig)
	if err != nil {
		log.Fatal("Error connecting to SFTP:", err)
	}
	defer sftpConnection.Close()

	sftpClient, err := sftp.NewClient(sftpConnection)
	if err != nil {
		log.Fatal("Error connecting to SFTP:", err)
	}
	defer sftpClient.Close()

	// List files in the root directory
	files, err := sftpClient.ReadDir("/")
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	fmt.Println("Files in the root directory:")
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
