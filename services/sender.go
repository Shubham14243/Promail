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

	Subject := "ProMail | Forgot your password? Let's fix that!"
	EmailType := "HTML"
	body := `
		<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Reset Password — ProMail</title>
</head>
<body>
<div style="background-color:#f4f5f7; padding:32px 12px;">
        <div style="max-width:560px; margin:0 auto; background-color:#171a21; border-radius:14px; overflow:hidden; border:1px solid #232733;">

          
          <div style="background:linear-gradient(135deg,#7c6cf5 0%,#5b8cff 100%); padding:28px 32px;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0">
              <tbody><tr>
                <td style="font-size:20px; font-weight:800; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">
                  ✉ ProMail
                </td>
                <td align="right" style="font-size:12px; color:rgba(255,255,255,.85); font-family:Helvetica,Arial,sans-serif;">
                  Security notice
                </td>
              </tr>
            </tbody></table>
          </div>

          
          <div style="padding:36px 32px;">
            <p style="margin:0 0 6px; font-size:13px; color:#7c6cf5; font-weight:700; font-family:Helvetica,Arial,sans-serif; letter-spacing:1px; text-transform:uppercase;">
              Password reset 🔑
            </p>
            <h1 style="margin:0 0 16px; font-size:26px; line-height:1.3; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">
              Forgot your password?
            </h1>
            <p style="margin:0 0 24px; font-size:15px; line-height:1.7; color:#aab0c0; font-family:Helvetica,Arial,sans-serif;">
              No worries — it happens. Click the button below to choose a new
              password. This link will expire in
              <strong style="color:#ffffff;">15 minutes</strong> for your security.
            </p>

            
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin:0 0 24px;">
              <tbody><tr>
                <td align="center">
                  <table role="presentation" cellpadding="0" cellspacing="0">
                    <tbody><tr>
                      <td style="border-radius:9px; background-color:#7c6cf5;">
                        <a href="{{reset_url}}" target="_blank" style="display:inline-block; padding:15px 42px; font-size:16px; font-weight:700; color:#ffffff; text-decoration:none; font-family:Helvetica,Arial,sans-serif;">
                          Reset my password
                        </a>
                      </td>
                    </tr>
                  </tbody></table>
                </td>
              </tr>
            </tbody></table>

            
            <div style="background-color:#1e222c; border-radius:10px; padding:16px 18px; margin:0 0 24px;">
              <p style="margin:0 0 6px; font-size:12px; color:#8b91a3; font-family:Helvetica,Arial,sans-serif;">
                Button not working? Paste this link into your browser:
              </p>
              <p style="margin:0; font-size:12px; word-break:break-all; font-family:ui-monospace,Menlo,Consolas,monospace;">
                <a href="{{reset_url}}" style="color:#7c6cf5; text-decoration:underline;">{{reset_url}}</a>
              </p>
            </div>

            
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0">
              <tbody><tr>
                <td style="border-left:3px solid #febc2e; padding:4px 0 4px 14px;">
                  <p style="margin:0; font-size:13px; line-height:1.6; color:#aab0c0; font-family:Helvetica,Arial,sans-serif;">
                    <strong style="color:#ffffff;">Didn't request this?</strong>
                    You can safely ignore this email — your password won't change
                    unless you click the link above.
                  </p>
                </td>
              </tr>
            </tbody></table>
          </div>

          
          <div style="padding:20px 32px; border-top:1px solid #232733; background-color:#14161c;">
            <p style="margin:0; font-size:12px; line-height:1.6; color:#6b7280; font-family:Helvetica,Arial,sans-serif;">
              This is an automated security message from ProMail. Please do not reply.
            </p>
            <p style="margin:10px 0 0; font-size:12px; color:#6b7280; font-family:Helvetica,Arial,sans-serif;">
              © {{year}} ProMail ·
            </p>
          </div>
        </div>
      </div>
</body>
</html>`
	body = strings.ReplaceAll(body, "{{reset_url}}", resetURL)
	body = strings.ReplaceAll(body, "{{year}}", strconv.Itoa(time.Now().Year()))
	IdemKey := uuid.NewString()

	if err := SendEmail(&fakeConfig, to, Subject, body, EmailType, IdemKey); err != nil {
		return err
	}

	return nil

}

func SendSignUpEmail(to string, name string) error {

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

	Subject := "Welcome to ProMail! Your SMTP is connected"
	EmailType := "HTML"
	body := `
		<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Welcome — ProMail</title>
</head>
<body>
<div style="background-color:#f4f5f7; padding:32px 12px;">
        <div style="max-width:560px; margin:0 auto; background-color:#171a21; border-radius:14px; overflow:hidden; border:1px solid #232733;">

          
          <div style="background:linear-gradient(135deg,#7c6cf5 0%,#5b8cff 100%); padding:28px 32px;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0">
              <tbody><tr>
                <td style="font-size:20px; font-weight:800; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">
                  ✉ ProMail
                </td>
                <td align="right" style="font-size:12px; color:rgba(255,255,255,.85); font-family:Helvetica,Arial,sans-serif;">
                  Free. No domain required.
                </td>
              </tr>
            </tbody></table>
          </div>

          
          <div style="padding:36px 32px;">
            <p style="margin:0 0 6px; font-size:13px; color:#7c6cf5; font-weight:700; font-family:Helvetica,Arial,sans-serif; letter-spacing:1px; text-transform:uppercase;">
              Welcome aboard 🎉
            </p>
            <h1 style="margin:0 0 16px; font-size:26px; line-height:1.3; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">
              Hi {{name}}, your SMTP is connected.
            </h1>
            <p style="margin:0 0 20px; font-size:15px; line-height:1.7; color:#aab0c0; font-family:Helvetica,Arial,sans-serif;">
              Thanks for joining ProMail. Your account is ready — start sending
              transactional email through the SMTP you already have. No domain to
              verify, no plan to pick, nothing to pay.
            </p>

            
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin:0 0 28px;">
              <tbody><tr>
                <td style="padding:16px; background-color:#1e222c; border-radius:10px; width:33%;" align="center">
                  <div style="font-size:22px; font-weight:800; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">$0</div>
                  <div style="font-size:11px; color:#8b91a3; font-family:Helvetica,Arial,sans-serif;">per month</div>
                </td>
                <td style="width:8px;">&nbsp;</td>
                <td style="padding:16px; background-color:#1e222c; border-radius:10px; width:33%;" align="center">
                  <div style="font-size:22px; font-weight:800; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">0</div>
                  <div style="font-size:11px; color:#8b91a3; font-family:Helvetica,Arial,sans-serif;">domains required</div>
                </td>
                <td style="width:8px;">&nbsp;</td>
                <td style="padding:16px; background-color:#1e222c; border-radius:10px; width:33%;" align="center">
                  <div style="font-size:22px; font-weight:800; color:#ffffff; font-family:Helvetica,Arial,sans-serif;">Many</div>
                  <div style="font-size:11px; color:#8b91a3; font-family:Helvetica,Arial,sans-serif;">SMTP configs at once</div>
                </td>
              </tr>
            </tbody></table>

            
            <table role="presentation" cellpadding="0" cellspacing="0">
              <tbody><tr>
                <td style="border-radius:9px; background-color:#7c6cf5;">
                  <a href="{{dashboard_url}}" target="_blank" style="display:inline-block; padding:14px 30px; font-size:15px; font-weight:700; color:#ffffff; text-decoration:none; font-family:Helvetica,Arial,sans-serif;">
                    Open your dashboard →
                  </a>
                </td>
                <td style="padding-left:16px;">
                  <a href="{{docs_url}}" style="font-size:14px; color:#aab0c0; text-decoration:underline; font-family:Helvetica,Arial,sans-serif;">
                    Read the docs
                  </a>
                </td>
              </tr>
            </tbody></table>
          </div>

          
          <div style="padding:20px 32px; border-top:1px solid #232733; background-color:#14161c;">
            <p style="margin:0; font-size:12px; line-height:1.6; color:#6b7280; font-family:Helvetica,Arial,sans-serif;">
              You're receiving this email because you created a ProMail account.
              If this wasn't you, you can safely ignore it.
            </p>
            <p style="margin:10px 0 0; font-size:12px; color:#6b7280; font-family:Helvetica,Arial,sans-serif;">
              © {{year}} ProMail ·
            </p>
          </div>
        </div>
      </div>
</body>
</html>`

	uiURL := strings.TrimRight(os.Getenv("UI_URL"), "/")

	body = strings.ReplaceAll(body, "{{name}}", name)
	body = strings.ReplaceAll(body, "{{dashboard_url}}", uiURL+"/home")
	body = strings.ReplaceAll(body, "{{docs_url}}", name+"/docs")
	body = strings.ReplaceAll(body, "{{year}}", strconv.Itoa(time.Now().Year()))
	IdemKey := uuid.NewString()

	if err := SendEmail(&fakeConfig, to, Subject, body, EmailType, IdemKey); err != nil {
		return err
	}

	return nil

}
