package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"lostmediago/internal/utils"
	"lostmediago/pkg/mq"

	"github.com/streadway/amqp"
)

// StartEmailWorker starts consuming email jobs from RabbitMQ
func StartEmailWorker() error {
	if mq.Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	msgs, err := mq.Channel.Consume(
		"email_queue", // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("[EMAIL WORKER] Received email message from queue")
			if err := processEmailMessage(msg); err != nil {
				log.Printf("[EMAIL WORKER ERROR] Error processing email message: %v", err)
				// Reject message and requeue it (only for transient errors)
				// For permanent errors, set requeue to false
				msg.Nack(false, true)
				continue
			}
			// Acknowledge message
			msg.Ack(false)
			log.Printf("[EMAIL WORKER] Email message processed successfully")
		}
	}()

	log.Printf("[EMAIL WORKER] Email worker started and listening for messages")

	return nil
}

func processEmailMessage(msg amqp.Delivery) error {
	var emailMsg mq.EmailMessage
	if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
		return fmt.Errorf("failed to unmarshal email message: %w", err)
	}

	switch emailMsg.Type {
	case "verification":
		log.Printf("[EMAIL WORKER] Processing verification email for: %s", emailMsg.To)
		err := utils.SendVerificationEmail(emailMsg.To, emailMsg.Token)
		if err != nil {
			log.Printf("[EMAIL WORKER ERROR] Failed to send verification email to %s: %v", emailMsg.To, err)
			return fmt.Errorf("failed to send verification email: %w", err)
		}
		log.Printf("[EMAIL WORKER SUCCESS] Verification email sent to: %s", emailMsg.To)

	case "reset_password":
		log.Printf("[EMAIL WORKER] Processing password reset email for: %s", emailMsg.To)
		err := utils.SendPasswordResetEmail(emailMsg.To, emailMsg.Token)
		if err != nil {
			log.Printf("[EMAIL WORKER ERROR] Failed to send password reset email to %s: %v", emailMsg.To, err)
			return fmt.Errorf("failed to send reset password email: %w", err)
		}
		log.Printf("[EMAIL WORKER SUCCESS] Password reset email sent to: %s", emailMsg.To)

	default:
		log.Printf("[EMAIL WORKER ERROR] Unknown email type: %s", emailMsg.Type)
		return fmt.Errorf("unknown email type: %s", emailMsg.Type)
	}

	return nil
}
