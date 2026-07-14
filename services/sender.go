package services

import (
	"net/smtp"
	"strconv"
)

func SendEmail(smtpHost string, smtpPort int, smtpUsername, smtpPassword, to, subject, body string) error {

	if err := validateEmailData(to); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	email_body := "From: Shubham Gupta\r\n" + "To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + body + "\r\n"

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, smtpUsername, []string{to}, []byte(email_body))
	if err != nil {
		return err
	}
	return nil
}
