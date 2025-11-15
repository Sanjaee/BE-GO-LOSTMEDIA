package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"lostmediago/internal/repositories"
	"lostmediago/pkg/mq"

	"github.com/streadway/amqp"
)

// StartUserActivityWorker starts consuming user activity events from RabbitMQ
func StartUserActivityWorker(userRepo repositories.UserRepository) error {
	if mq.Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	msgs, err := mq.Channel.Consume(
		"user_activity_queue", // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("[USER ACTIVITY WORKER] Received user activity event from queue")
			if err := processUserActivityMessage(msg, userRepo); err != nil {
				log.Printf("[USER ACTIVITY WORKER ERROR] Error processing user activity event: %v", err)
				// Reject message and requeue it (only for transient errors)
				// For permanent errors, set requeue to false
				msg.Nack(false, true)
				continue
			}
			// Acknowledge message
			msg.Ack(false)
			log.Printf("[USER ACTIVITY WORKER] User activity event processed successfully")
		}
	}()

	log.Printf("[USER ACTIVITY WORKER] User activity worker started and listening for messages")

	return nil
}

func processUserActivityMessage(msg amqp.Delivery, userRepo repositories.UserRepository) error {
	var event mq.UserActivityEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user activity event: %w", err)
	}

	log.Printf("[USER ACTIVITY WORKER] Processing event: type=%s, user_id=%s, email=%s", event.Type, event.UserID, event.Email)

	switch event.Type {
	case "login":
		// Update last login timestamp
		if err := userRepo.UpdateLastLogin(event.UserID); err != nil {
			log.Printf("[USER ACTIVITY WORKER ERROR] Failed to update last login for user %s: %v", event.UserID, err)
			return fmt.Errorf("failed to update last login: %w", err)
		}
		log.Printf("[USER ACTIVITY WORKER SUCCESS] Updated last login for user: %s", event.UserID)

		// Here you can add more logic:
		// - Send analytics event
		// - Update user activity metrics
		// - Trigger notifications
		// - Log security events

	case "register":
		log.Printf("[USER ACTIVITY WORKER SUCCESS] User registered: %s (%s)", event.UserID, event.Email)

		// Here you can add more logic:
		// - Send welcome email (if not sent during registration)
		// - Add to analytics
		// - Create default settings
		// - Trigger onboarding workflow

	case "email_verified":
		log.Printf("[USER ACTIVITY WORKER SUCCESS] Email verified for user: %s", event.UserID)

		// Here you can add more logic:
		// - Send confirmation email
		// - Update user metrics
		// - Trigger post-verification workflows

	case "password_reset":
		log.Printf("[USER ACTIVITY WORKER SUCCESS] Password reset for user: %s", event.UserID)

		// Here you can add more logic:
		// - Send security notification
		// - Log security event
		// - Update user metrics

	default:
		log.Printf("[USER ACTIVITY WORKER ERROR] Unknown event type: %s", event.Type)
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	return nil
}
