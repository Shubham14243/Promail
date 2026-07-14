package services

import (
	"errors"
	"promail/models"
	"regexp"
)

// EMAIL Validation
var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// NAME Validation
var nameRegex = regexp.MustCompile(
	`^[a-zA-Z ]{2,50}$`,
)

func IsValidName(name string) bool {
	return nameRegex.MatchString(name)
}

// APPNAME Validation
var appNameRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9 ._/-]{4,49}$`,
)

func IsValidAppName(appName string) bool {
	return appNameRegex.MatchString(appName)
}

// SUBJECT Validation
var subjectRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9 ._/-|]{4,99}$`,
)

func IsValidSubject(subject string) bool {
	return subjectRegex.MatchString(subject)
}

// PASSWORD Validation (Go-safe version)
// NOTE: No lookaheads allowed in Go regexp
var (
	lowerRegex  = regexp.MustCompile(`[a-z]`)
	upperRegex  = regexp.MustCompile(`[A-Z]`)
	numberRegex = regexp.MustCompile(`[0-9]`)
)

// SLUG Validation
var slugRegex = regexp.MustCompile(
	`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
)

func IsValidSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}

func IsValidPassword(password string) bool {

	if len(password) < 8 || len(password) > 25 {
		return false
	}

	if !lowerRegex.MatchString(password) {
		return false
	}

	if !upperRegex.MatchString(password) {
		return false
	}

	if !numberRegex.MatchString(password) {
		return false
	}

	return true
}

// SIGNUP VALIDATION
func ValidateUserCreate(req models.UserCreateRequest) error {

	if !IsValidName(req.Name) {
		return errors.New("Invalid name: 2-50 chars, letters and spaces only.")
	}

	if !IsValidEmail(req.Email) {
		return errors.New("Invalid email.")
	}

	if !IsValidPassword(req.Password) {
		return errors.New("Invalid password: must be 8-25 chars with uppercase, lowercase and number.")
	}

	return nil
}

// LOGIN VALIDATION
func ValidateLogin(req models.LoginRequest) error {

	if !IsValidEmail(req.Email) {
		return errors.New("Invalid email.")
	}

	if !IsValidPassword(req.Password) {
		return errors.New("Invalid password: must be 8-25 chars with uppercase, lowercase and number.")
	}

	return nil
}

// USER UPDATE VALIDATION
func ValidateUserUpdate(req models.UserUpdateRequest) error {

	if !IsValidName(req.Name) {
		return errors.New("Invalid name: 2-50 chars, letters and spaces only.")
	}

	if !IsValidEmail(req.Email) {
		return errors.New("Invalid email.")
	}

	return nil
}

// APP VALIDATION
func ValidateAppCreate(req models.CreateApp) error {

	if !IsValidAppName(req.Name) {
		return errors.New("Invalid app name: 5-50 chars, starting with letter, letters, number, spaces and characters['.','_','-','/'] only.")
	}
	return nil
}

func ValidateAppUpdate(req models.UpdateApp) error {

	if !IsValidAppName(req.Name) {
		return errors.New("Invalid app name: 5-50 chars, starting with letter, letters, number, spaces and characters['.','_','-','/'] only.")
	}

	if !IsValidName(req.Status) {
		return errors.New("Invalid status: must be either 'active' or 'inactive'.")
	}

	return nil
}

func ValidateTemplateCreate(req models.TemplateCreate) error {

	if !IsValidAppName(req.Name) {
		return errors.New("Invalid template name: 5-50 chars, starting with letter, letters, number, spaces and characters['.','_','-','/'] only.")
	}

	if !IsValidSlug(req.Slug) {
		return errors.New("Invalid template slug: 5-50 chars, lowercase letters, numbers and hyphens only.")
	}

	if !IsValidSubject(req.Subject) {
		return errors.New("Invalid template subject: 5-100 chars, letters, numbers and spaces only.")
	}

	if req.Type != "html" && req.Type != "text" {
		return errors.New("Invalid template type: must be either 'html' or 'text'.")
	}

	if req.Content == "" {
		return errors.New("Invalid template content: cannot be empty.")
	}

	return nil
}

func ValidateTemplateUpdate(req models.TemplateUpdate) error {

	if !IsValidAppName(req.Name) {
		return errors.New("Invalid template name: 5-50 chars, starting with letter, letters, number, spaces and characters['.','_','-','/'] only.")
	}

	if !IsValidSlug(req.Slug) {
		return errors.New("Invalid template slug: 5-50 chars, lowercase letters, numbers and hyphens only.")
	}

	if !IsValidSubject(req.Subject) {
		return errors.New("Invalid template subject: 5-100 chars, letters, numbers and spaces only.")
	}

	if req.Status != "active" && req.Status != "inactive" {
		return errors.New("Invalid template status: must be either 'active' or 'inactive'.")
	}

	return nil
}

func ValidateTemplateContent(req models.TemplateContent) error {

	if req.Type != "html" && req.Type != "text" {
		return errors.New("Invalid template type: must be either 'html' or 'text'.")
	}

	if req.Content == "" {
		return errors.New("Invalid template content: cannot be empty.")
	}

	return nil
}

// APP CONFIG VALIDATION
func ValidateAppConfigCreate(req models.AppConfigCreate) error {

	if req.AppID <= 0 {
		return errors.New("Invalid app_id: must be a positive integer.")
	}

	if req.SMTPHost == "" {
		return errors.New("Invalid smtp_host: cannot be empty.")
	}

	if req.SMTPPort <= 0 || req.SMTPPort > 65535 {
		return errors.New("Invalid smtp_port: must be a positive integer between 1 and 65535.")
	}

	if !IsValidName(req.SMTPName) {
		return errors.New("Invalid smtp_name: must be a valid name.")
	}

	if !IsValidEmail(req.SMTPUsername) {
		return errors.New("Invalid smtp_username: must be a valid email address.")
	}

	if req.SMTPPassword == "" {
		return errors.New("Invalid smtp_password: cannot be empty.")
	}

	if req.OpenTrack != "active" && req.OpenTrack != "inactive" {
		return errors.New("Invalid open_track: must be either 'active' or 'inactive'.")
	}

	if req.ClickTrack != "active" && req.ClickTrack != "inactive" {
		return errors.New("Invalid click_track: must be either 'active' or 'inactive'.")
	}

	if req.AutoRetry != "active" && req.AutoRetry != "inactive" {
		return errors.New("Invalid auto_retry: must be either 'active' or 'inactive'.")
	}

	if req.RetryMaxCount < 0 {
		return errors.New("Invalid retry_max_count: must be a non-negative integer.")
	}

	return nil
}

func ValidateAppConfigUpdate(req models.AppConfigUpdate) error {

	if req.SMTPHost == "" {
		return errors.New("Invalid smtp_host: cannot be empty.")
	}

	if req.SMTPPort <= 0 || req.SMTPPort > 65535 {
		return errors.New("Invalid smtp_port: must be a positive integer between 1 and 65535.")
	}

	if !IsValidName(req.SMTPName) {
		return errors.New("Invalid smtp_name: must be a valid name.")
	}

	if !IsValidEmail(req.SMTPUsername) {
		return errors.New("Invalid smtp_username: must be a valid email address.")
	}

	if req.SMTPPassword == "" {
		return errors.New("Invalid smtp_password: cannot be empty.")
	}

	if req.OpenTrack != "active" && req.OpenTrack != "inactive" {
		return errors.New("Invalid open_track: must be either 'active' or 'inactive'.")
	}

	if req.ClickTrack != "active" && req.ClickTrack != "inactive" {
		return errors.New("Invalid click_track: must be either 'active' or 'inactive'.")
	}

	if req.AutoRetry != "active" && req.AutoRetry != "inactive" {
		return errors.New("Invalid auto_retry: must be either 'active' or 'inactive'.")
	}

	if req.RetryMaxCount < 0 {
		return errors.New("Invalid retry_max_count: must be a non-negative integer.")
	}

	return nil
}

// Send EMAIL Validation

func validateEmailData(receiver string) error {

	if !IsValidEmail(receiver) {
		return errors.New("Invalid receiver email.")
	}

	return nil
}
