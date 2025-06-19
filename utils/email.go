package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	// Check if we're in development mode (no email credentials)
	emailFrom := os.Getenv("EMAIL_FROM")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	// If credentials are placeholders or empty, just log instead of sending
	if emailFrom == "" || emailFrom == "your-email@gmail.com" ||
		emailPassword == "" || emailPassword == "your-email-password" {
		log.Println("==== DEVELOPMENT MODE: EMAIL NOT SENT ====")
		log.Printf("To: %s\n", to)
		log.Printf("Subject: %s\n", subject)
		log.Printf("Body: %s\n", body)
		log.Println("==========================================")
		return nil // Return success in dev mode
	}

	// Otherwise proceed with actual email sending
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", emailFrom, to, subject, body)

	auth := smtp.PlainAuth("", emailFrom, emailPassword, smtpHost)
	addr := smtpHost + ":" + smtpPort

	err := smtp.SendMail(addr, auth, emailFrom, []string{to}, []byte(msg))
	return err
}
