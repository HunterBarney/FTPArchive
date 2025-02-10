package emailclient

import (
	"FTPArchive/internal/config"
	"fmt"
	"log"
)

func SendEmail(subject, body string, config *config.Config, profile *config.Profile) error {
	// TODO: TESTING ONLY, REMOVE BEFORE FINAL VERSION OF BRANCH
	log.Printf("Config { Host %s, port %s, username %s, password %s, from %s, to %s, cc %s, bcc %s}", config.SMTP.Host, config.SMTP.Port, config.SMTP.Username, config.SMTP.Password, config.SMTP.From, config.SMTP.To, config.SMTP.CC, config.SMTP.BCC)
	log.Printf("Profile { Host %s, port %s, username %s, password %s, from %s, to %s, cc %s, bcc %s}", profile.SMTP.Host, profile.SMTP.Port, profile.SMTP.Username, profile.SMTP.Password, profile.SMTP.From, profile.SMTP.To, profile.SMTP.CC, profile.SMTP.BCC)

	smtp, err := validateSMTP(profile, config)
	if err != nil {
		return err
	}

	// TODO: TESTING ONLY, REMOVE BEFORE FINAL VERSION OF BRANCH
	log.Printf("Combined { Host %s, port %s, username %s, password %s, from %s, to %s, cc %s, bcc %s}", smtp.Host, smtp.Port, smtp.Username, smtp.Password, smtp.From, smtp.To, smtp.CC, smtp.BCC)

	return nil
}

func validateSMTP(profile *config.Profile, cfg *config.Config) (config.SMTPInfo, error) {
	smtp := config.SMTPInfo{
		Host:     getFirstNonBlank(profile.SMTP.Host, cfg.SMTP.Host),
		Port:     getFirstNonZero(profile.SMTP.Port, cfg.SMTP.Port),
		Username: getFirstNonBlank(profile.SMTP.Username, cfg.SMTP.Username),
		Password: getFirstNonBlank(profile.SMTP.Password, cfg.SMTP.Password),
		From:     getFirstNonBlank(profile.SMTP.From, cfg.SMTP.From),
		To:       getFirstNonEmpty(profile.SMTP.To, cfg.SMTP.To),
		CC:       getFirstNonEmpty(profile.SMTP.CC, cfg.SMTP.CC),
		BCC:      getFirstNonEmpty(profile.SMTP.BCC, cfg.SMTP.BCC),
	}

	if smtp.Host == "" {
		return smtp, fmt.Errorf("SMTP Host is empty in profile AND config")
	}

	if smtp.Port == 0 {
		return smtp, fmt.Errorf("SMTP Port is empty in profile AND config")
	}

	if smtp.Username == "" {
		return smtp, fmt.Errorf("SMTP Username is empty in profile AND config")
	}

	if smtp.Password == "" {
		return smtp, fmt.Errorf("SMTP Password is empty in profile AND config")
	}

	if smtp.From == "" {
		return smtp, fmt.Errorf("SMTP From is empty in profile AND config")
	}

	if smtp.To == nil && smtp.CC == nil && smtp.BCC == nil {
		return smtp, fmt.Errorf("SMPT to, cc, and BCC are empty in profile AND config")
	}

	if len(smtp.To) == 0 && len(smtp.CC) == 0 && len(smtp.BCC) == 0 {
		return smtp, fmt.Errorf("SMPT to, cc, and BCC are empty in profile AND config")
	}
	return smtp, nil
}

// getFirstNonBlank returns the first argument if it is a non-empty string, otherwise it returns b regardless if its empty
func getFirstNonBlank(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// getFirstNonEmpty returns a if it is a non-empty string array, otherwise returns b regardless if it is empty
func getFirstNonEmpty(a, b []string) []string {
	if len(a) > 0 {
		return a
	}
	return b
}

// getFirstNonZero returns the first argument if it is not zero, otherwise it returns b regardless if its zero
func getFirstNonZero(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}
