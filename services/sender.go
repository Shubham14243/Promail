package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"promail/models"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SendEmail(appConf *models.AppConfigData, to string, subject string, body string, emailType string, idemKey string) error {

	relayURL := os.Getenv("RELAY_URL")
	if relayURL == "" {
		relayURL = "http://127.0.0.1:8090/relay"
	}
	apiKey := os.Getenv("RELAY_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RELAY_API_KEY is not configured")
	}

	payload := map[string]any{
		"host":     appConf.SMTPHost,
		"port":     appConf.SMTPPort,
		"username": appConf.SMTPUsername,
		"password": appConf.SMTPPassword,
		"from":     appConf.SMTPName,
		"to":       to,
		"subject":  subject,
		"body":     body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal relay payload: %w", err)
	}

	encrypted, err := Encrypt(string(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to encrypt relay payload: %w", err)
	}

	wrapped := map[string]string{"payload": encrypted}
	wrappedBytes, err := json.Marshal(wrapped)
	if err != nil {
		return fmt.Errorf("failed to marshal wrapped payload: %w", err)
	}

	if idemKey == "" {
		idemKey = uuid.NewString()
	}
	req, err := http.NewRequest(http.MethodPost, relayURL, bytes.NewReader(wrappedBytes))
	if err != nil {
		return fmt.Errorf("failed to build relay request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Idempotency-Key", idemKey)

	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("relay request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var relayErr struct {
		Error  string         `json:"error"`
		Result map[string]any `json:"result"`
	}
	if decErr := json.NewDecoder(resp.Body).Decode(&relayErr); decErr != nil {
		return fmt.Errorf("relay returned status %d", resp.StatusCode)
	}
	if relayErr.Result != nil {
		if msg, ok := relayErr.Result["message"].(string); ok && msg != "" {
			return fmt.Errorf("relay status %d: %s", resp.StatusCode, msg)
		}
	}
	if relayErr.Error != "" {
		return fmt.Errorf("relay status %d: %s", resp.StatusCode, relayErr.Error)
	}
	return fmt.Errorf("relay returned status %d", resp.StatusCode)
}

func PrepareEmailBody(body string, variables map[string]string) string {

	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		body = strings.ReplaceAll(body, placeholder, value)
	}

	return body
}

func AddOpenTrackingBody(body string, openUUID string, tempType string) string {

	baseUrl := os.Getenv("APP_BASE_URL")
	openStr := "<img src='" + baseUrl + "/api/v1/email/track/open/" + openUUID + "' style='display:none!important;max-width:1px;max-height:1px;'/>"

	if tempType == "text" {
		body += openStr
	} else {
		openStr = "<img src='" + baseUrl + "/api/v1/email/track/open/" + openUUID + "' style='display:none!important;max-width:1px;max-height:1px;'/>" + "</body>"
		body = strings.ReplaceAll(body, "</body>", openStr)
	}

	return body
}

func AddClickTrackingBody(body string) (string, []models.ClickTracking) {
	baseURL := os.Getenv("APP_BASE_URL")

	re := regexp.MustCompile(`(?i)href\s*=\s*("([^"]*)"|'([^']*)')`)

	var trackings []models.ClickTracking

	updated := re.ReplaceAllStringFunc(body, func(match string) string {
		sub := re.FindStringSubmatch(match)

		var originalURL string
		if sub[2] != "" {
			originalURL = sub[2]
		} else {
			originalURL = sub[3]
		}

		parsedURL, err := url.Parse(originalURL)
		if err != nil || parsedURL.Scheme != "" && parsedURL.Scheme != "http" && parsedURL.Scheme != "https" || parsedURL.Scheme == "" && parsedURL.Host != "" {
			return match
		}

		token := uuid.New()

		trackings = append(trackings, models.ClickTracking{
			Token:       token,
			OriginalURL: originalURL,
		})

		return fmt.Sprintf(`href="%s/api/v1/email/track/click/%s"`,
			baseURL,
			token.String(),
		)
	})

	return updated, trackings
}

func SendResetEmail(to string, resetURL string) error {

	port, err := strconv.Atoi(os.Getenv("RESET_PORT"))
	if err != nil {
		port = 587 // default
	}

	fakeConfig := models.AppConfigData{
		ID:            0,
		AppID:         0,
		SMTPHost:      os.Getenv("RESET_HOST"),
		SMTPPort:      port,
		SMTPName:      "ProMail",
		SMTPUsername:  os.Getenv("RESET_EMAIL"),
		SMTPPassword:  os.Getenv("RESET_PASSWORD"),
		OpenTrack:     "",
		ClickTrack:    "",
		AutoRetry:     "",
		RetryMaxCount: 0,
		CreatedAt:     "",
		UpdatedAt:     "",
	}

	fmt.Print(resetURL)

	Subject := "ProMail | Reset Your Password!"
	EmailType := "HTML"
	body := `
		<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset your password</title>
</head>

<body style="
    margin: 0;
    padding: 0;
    background-color: #ffffff;
    font-family: Arial, Helvetica, sans-serif;
    color: #111111;
">
<div></div>
    <table
        role="presentation"
        width="100%"
        cellspacing="0"
        cellpadding="0"
        border="0"
        style="background-color: #ffffff;"
    >
        <tr>
            <td align="center" style="padding: 18px 20px;">

                <!-- Main Card -->
                <table
                    role="presentation"
                    width="100%"
                    cellspacing="0"
                    cellpadding="0"
                    border="0"
                    style="
                        max-width: 640px;
                        border: 1px solid #dddddd;
                        border-radius: 22px;
                        background-color: #ffffff;
                    "
                >
                    <tr>
                        <td style="padding: 38px 40px 36px 40px;">

                            <!-- Heading -->
                            <h1 style="
                                margin: 0 0 14px 0;
                                font-size: 22px;
                                line-height: 30px;
                                font-weight: 600;
                                color: #080808;
                            ">
                                Reset your password
                            </h1>

                            <!-- Description -->
                            <p style="
                                margin: 0 0 24px 0;
                                font-size: 16px;
                                line-height: 24px;
                                color: #111111;
                            ">
                                We received a request to reset the password for your account.
                            </p>

                            <!-- Button -->
                            <table
                                role="presentation"
                                cellspacing="0"
                                cellpadding="0"
                                border="0"
                            >
                                <tr>
                                    <td
                                        style="
                                            border-radius: 6px;
                                            background-color: #000000;
                                        "
                                    >
                                        <a
                                            href="{{RESET_URL}}"
                                            target="_blank"
                                            style="
                                                display: inline-block;
                                                padding: 9px 14px;
                                                font-size: 14px;
                                                line-height: 18px;
                                                font-weight: 600;
                                                color: #ffffff;
                                                text-decoration: none;
                                                border-radius: 6px;
                                            "
                                        >
                                            Reset password
                                        </a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Divider -->
                            <table
                                role="presentation"
                                width="100%"
                                cellspacing="0"
                                cellpadding="0"
                                border="0"
                                style="margin-top: 32px;"
                            >
                                <tr>
                                    <td style="
                                        border-top: 1px solid #dddddd;
                                        font-size: 0;
                                        line-height: 0;
                                    ">
                                        &nbsp;
                                    </td>
                                </tr>
                            </table>

                            <!-- Expiry -->
                            <p style="
                                margin: 30px 0 10px 0;
                                font-size: 13px;
                                line-height: 20px;
                                color: #666666;
                            ">
                                This link expires in 15 minutes.
                            </p>

                            <!-- Security Notice -->
                            <p style="
                                margin: 0;
                                font-size: 13px;
                                line-height: 20px;
                                color: #666666;
                            ">
                                If you didn’t request to reset your password, you can safely
                                ignore this email. Someone else might have typed your email
                                address by mistake.
                            </p>

                        </td>
                    </tr>
                </table>

            </td>
        </tr>
    </table>

</body>
</html>`
	body = strings.ReplaceAll(body, "{{RESET_URL}}", resetURL)
	IdemKey := uuid.NewString()

	if err := SendEmail(&fakeConfig, to, Subject, body, EmailType, IdemKey); err != nil {
		return err
	}

	return nil

}
