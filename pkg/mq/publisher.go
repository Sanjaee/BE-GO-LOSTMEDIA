package mq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// EmailMessage represents an email job message
type EmailMessage struct {
	Type    string `json:"type"`    // "verification" or "reset_password"
	To      string `json:"to"`      // recipient email
	Token   string `json:"token"`   // verification or reset token
	Subject string `json:"subject"` // email subject
	Body    string `json:"body"`    // email body
}

// PublishEmail publishes an email job to the email queue
func PublishEmail(message *EmailMessage) error {
	if Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	err = Channel.Publish(
		"",            // exchange
		"email_queue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}

	return nil
}

// PublishVerificationEmail publishes an email verification job
func PublishVerificationEmail(to, token string) error {
	message := &EmailMessage{
		Type:    "verification",
		To:      to,
		Token:   token,
		Subject: "Verify Your Email - LostMediaGo",
	}

	return PublishEmail(message)
}

// PublishPasswordResetEmail publishes a password reset email job
func PublishPasswordResetEmail(to, token string) error {
	message := &EmailMessage{
		Type:    "reset_password",
		To:      to,
		Token:   token,
		Subject: "Reset Your Password - LostMediaGo",
	}

	return PublishEmail(message)
}
