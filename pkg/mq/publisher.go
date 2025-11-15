package mq

import (
	"encoding/json"
	"fmt"
	"time"

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

// UserActivityEvent represents user activity events
type UserActivityEvent struct {
	Type      string                 `json:"type"`               // "login", "logout", "register", "password_reset", "email_verified"
	UserID    string                 `json:"user_id"`            // user ID
	Email     string                 `json:"email"`              // user email
	Timestamp time.Time              `json:"timestamp"`          // event timestamp
	Metadata  map[string]interface{} `json:"metadata,omitempty"` // additional metadata
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

// PublishUserActivity publishes a user activity event
func PublishUserActivity(event *UserActivityEvent) error {
	if Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user activity event: %w", err)
	}

	err = Channel.Publish(
		"",                    // exchange
		"user_activity_queue", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish user activity event: %w", err)
	}

	return nil
}

// PublishLoginEvent publishes a user login event
func PublishLoginEvent(userID, email string) error {
	event := &UserActivityEvent{
		Type:      "login",
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
	return PublishUserActivity(event)
}

// PublishRegisterEvent publishes a user registration event
func PublishRegisterEvent(userID, email string) error {
	event := &UserActivityEvent{
		Type:      "register",
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
	return PublishUserActivity(event)
}

// PublishEmailVerifiedEvent publishes an email verified event
func PublishEmailVerifiedEvent(userID, email string) error {
	event := &UserActivityEvent{
		Type:      "email_verified",
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
	return PublishUserActivity(event)
}

// PublishPasswordResetEvent publishes a password reset event
func PublishPasswordResetEvent(userID, email string) error {
	event := &UserActivityEvent{
		Type:      "password_reset",
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
	return PublishUserActivity(event)
}
