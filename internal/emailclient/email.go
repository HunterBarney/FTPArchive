package emailclient

import (
	"FTPArchive/internal/config"
	"fmt"
	"github.com/wneessen/go-mail"
	"os"
)

func SendEmail(subject, body string, config *config.Config, profile *config.Profile, logFile *os.File) error {
	smtp, err := validateSMTP(profile, config)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	err = msg.From(smtp.From)
	if err != nil {
		return err
	}

	for _, element := range smtp.To {
		err = msg.AddTo(element)
		if err != nil {
			return err
		}
	}

	for _, element := range smtp.CC {
		err = msg.AddCc(element)
		if err != nil {
			return err
		}
	}

	for _, element := range smtp.BCC {
		err = msg.AddBcc(element)
		if err != nil {
			return err
		}
	}

	if config.SendLogOverEmail {
		msg.AttachFile(logFile.Name())
	}

	client, err := mail.NewClient(smtp.Host, mail.WithPort(smtp.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtp.Username), mail.WithPassword(smtp.Password))
	if err != nil {
		return err
	}

	err = client.DialAndSend(msg)
	if err != nil {
		return err
	}

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

// SendFailEmail sends a pre-made fail email with body text of "Download {output name} failed with {error}".
// Automatically checks if sending emails is enabled in the config.
func SendFailEmail(config *config.Config, profile *config.Profile, e error, log *os.File) error {
	if config.SendLogOverEmail {
		body := fmt.Sprintf("Download %s failed with error %s", profile.OutputName, e)
		err := SendEmail("FTPArchive Failed!", body, config, profile, log)
		if err != nil {
			return err
		}
	}
	return nil
}
