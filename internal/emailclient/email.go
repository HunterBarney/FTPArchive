package emailclient

import (
	"FTPArchive/internal/config"
	"log"
)

func SendEmail(subject, body string, config *config.Config, profile *config.Profile) error {
	log.Printf("Config { Host %s, port %s, username %s, password %s, from %s, to %s, cc %s, bcc %s}", config.SMTP.Host, config.SMTP.Port, config.SMTP.Username, config.SMTP.Password, config.SMTP.From, config.SMTP.To, config.SMTP.CC, config.SMTP.BCC)
	log.Printf("Profil { Host %s, port %s, username %s, password %s, from %s, to %s, cc %s, bcc %s}", profile.SMTP.Host, profile.SMTP.Port, profile.SMTP.Username, profile.SMTP.Password, profile.SMTP.From, profile.SMTP.To, profile.SMTP.CC, profile.SMTP.BCC)

	return nil
}
