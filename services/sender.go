package services

import (
	"net/smtp"
	"os"
	"promail/models"
	"strconv"
	"strings"
)

func SendEmail(appConf *models.AppConfigData, to string, subject string, body string, emailType string) error {

	auth := smtp.PlainAuth("", appConf.SMTPUsername, appConf.SMTPPassword, appConf.SMTPHost)

	emailBody := ""
	emailBody += "From: " + appConf.SMTPName + " <" + appConf.SMTPUsername + ">\r\n"
	emailBody += "To: " + to + "\r\n"
	emailBody += "Subject: " + subject + "\r\n"
	emailBody += "MIME-Version: 1.0\r\n"
	emailBody += "Content-Type: text/html; charset=UTF-8\r\n"
	emailBody += "\r\n"
	emailBody += body

	err := smtp.SendMail(appConf.SMTPHost+":"+strconv.Itoa(appConf.SMTPPort), auth, appConf.SMTPUsername, []string{to}, []byte(emailBody))
	if err != nil {
		return err
	}
	return nil
}

func PrepareEmailBody(body string, variables map[string]string) string {

	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		body = strings.ReplaceAll(body, placeholder, value)
	}

	return body
}

func AddOpenTracking(body string, openUUID string, tempType string) string {

	baseUrl := os.Getenv("APP_BASE_URL")
	openStr := "<img src='" + baseUrl + "/api/v1/email/track/open/" + openUUID + "'/>"

	if tempType == "text" {
		body += openStr
	} else {
		openStr := "<img src='" + baseUrl + "/api/v1/email/track/open/" + openUUID + "'/>" + "</body>"
		body = strings.ReplaceAll(body, "</body>", openStr)
	}

	return body
}
