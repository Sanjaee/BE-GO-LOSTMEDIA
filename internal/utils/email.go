package utils

import (
	"fmt"
	"log"
	"lostmediago/internal/config"
	"net/smtp"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	cfg := config.AppConfig.Email

	if cfg.Username == "" || cfg.Password == "" {
		passwordStatus := "(empty)"
		if cfg.Password != "" {
			passwordStatus = "(set)"
		}
		log.Printf("[EMAIL ERROR] SMTP credentials not configured. Username: %s, Password: %s",
			cfg.Username, passwordStatus)
		log.Printf("[EMAIL] Email prepared (not sent) for: %s | Subject: %s", to, subject)
		return fmt.Errorf("SMTP credentials not configured")
	}

	from := cfg.From

	// Create HTML email with proper headers
	headers := "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += fmt.Sprintf("From: %s\r\n", from)
	headers += fmt.Sprintf("To: %s\r\n", to)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "\r\n"

	msg := []byte(headers + body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)

	log.Printf("[EMAIL] Attempting to send email to: %s | Subject: %s | SMTP: %s", to, subject, addr)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to send email to %s | Subject: %s | Error: %v", to, subject, err)
		log.Printf("[EMAIL ERROR] SMTP Config - Host: %s, Port: %d, From: %s, Username: %s",
			cfg.SMTPHost, cfg.SMTPPort, from, cfg.Username)
		return fmt.Errorf("failed to send email via SMTP to %s: %w", to, err)
	}

	log.Printf("[EMAIL SUCCESS] Email sent successfully to: %s | Subject: %s", to, subject)
	return nil
}

// SendVerificationEmail sends an email verification email with OTP code
func SendVerificationEmail(email, otpCode string) error {
	subject := "Verify Your Email - LostMedia"
	body := GetVerificationEmailTemplate(otpCode)

	log.Printf("[EMAIL] Sending verification email to: %s | OTP: %s", email, otpCode)

	err := SendEmail(email, subject, body)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to send verification email to %s: %v", email, err)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("[EMAIL SUCCESS] Verification email sent to: %s", email)
	return nil
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(email, token string) error {
	cfg := config.AppConfig
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", cfg.Server.FrontendURL, token)

	subject := "Reset Your Password - LostMedia"
	body := GetPasswordResetEmailTemplate(resetURL)

	log.Printf("[EMAIL] Sending password reset email to: %s", email)

	err := SendEmail(email, subject, body)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to send password reset email to %s: %v", email, err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Printf("[EMAIL SUCCESS] Password reset email sent to: %s", email)
	return nil
}
