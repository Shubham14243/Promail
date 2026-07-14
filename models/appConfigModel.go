package models

type AppConfigData struct {
	ID            int64  `json:"id"`
	AppID         int64  `json:"app_id"`
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPName      string `json:"smtp_name"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	OpenTrack     string `json:"open_track"`
	ClickTrack    string `json:"click_track"`
	AutoRetry     string `json:"auto_retry"`
	RetryMaxCount int    `json:"retry_max_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type AppConfigCreate struct {
	AppID         int64  `json:"app_id"`
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPName      string `json:"smtp_name"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	OpenTrack     string `json:"open_track"`
	ClickTrack    string `json:"click_track"`
	AutoRetry     string `json:"auto_retry"`
	RetryMaxCount int    `json:"retry_max_count"`
}

type AppConfigUpdate struct {
	ID            int64  `json:"id"`
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPName      string `json:"smtp_name"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	OpenTrack     string `json:"open_track"`
	ClickTrack    string `json:"click_track"`
	AutoRetry     string `json:"auto_retry"`
	RetryMaxCount int    `json:"retry_max_count"`
}
