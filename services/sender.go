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
