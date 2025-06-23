package services

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// EmailServiceProvider defines the interface for sending emails.
type EmailServiceProvider interface {
	SendEmail(to, subject, body string) error
}

// GmailService is an implementation of EmailServiceProvider for sending emails via Gmail.
type GmailService struct {
	dialer *gomail.Dialer
}

// NewEmailService creates a new email service provider based on the environment configuration.
func NewEmailService() (EmailServiceProvider, error) {
	provider := os.Getenv("EMAIL_PROVIDER")
	switch provider {
	case "gmail":
		email := os.Getenv("GMAIL_SENDER_EMAIL")
		password := os.Getenv("GMAIL_APP_PASSWORD")
		if email == "" || password == "" {
			return nil, fmt.Errorf("GMAIL_SENDER_EMAIL and GMAIL_APP_PASSWORD must be set for gmail provider")
		}
		// Gmail's SMTP server uses port 587
		d := gomail.NewDialer("smtp.gmail.com", 587, email, password)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // This is for local development, should be more secure in production
		return &GmailService{dialer: d}, nil
	// Add cases for other providers like "mailgun", "mailtrap" here.
	default:
		// A "noop" or "log" provider is useful for development or testing
		// where you don't want to send real emails.
		fmt.Printf("Email provider '%s' not configured, using log provider.\n", provider)
		return &LogEmailService{}, nil
	}
}

// SendEmail sends an email using the Gmail SMTP server.
func (s *GmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.dialer.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

// LogEmailService is an implementation of EmailServiceProvider that logs emails instead of sending them.
// Useful for development environments.
type LogEmailService struct{}

// SendEmail logs the email details to the console.
func (s *LogEmailService) SendEmail(to, subject, body string) error {
	fmt.Printf("\n--- New Email --- \n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("--- End Email --- \n\n")
	return nil
}
