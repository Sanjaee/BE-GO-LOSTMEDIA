package mq

import (
	"fmt"
	"log"
	"time"

	"lostmediago/internal/config"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// Connect initializes RabbitMQ connection with polling retry
func Connect() error {
	cfg := config.AppConfig.RabbitMQ

	log.Printf("[RABBITMQ] Connecting to RabbitMQ...")
	log.Printf("[RABBITMQ] Host: %s, Port: %d, User: %s, VHost: %s",
		cfg.Host, cfg.Port, cfg.User, cfg.VHost)

	var dsn string
	if cfg.VHost != "" && cfg.VHost != "/" {
		dsn = fmt.Sprintf("amqp://%s:%s@%s:%d%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.VHost,
		)
	} else {
		dsn = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
		)
	}

	// Poll until RabbitMQ is ready
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("[RABBITMQ] Attempting connection (attempt %d/%d)...", i+1, maxRetries)

		var err error
		Conn, err = amqp.Dial(dsn)
		if err == nil {
			log.Printf("[RABBITMQ] Connection established successfully")

			log.Printf("[RABBITMQ] Opening channel...")
			Channel, err = Conn.Channel()
			if err == nil {
				log.Printf("[RABBITMQ] Channel opened successfully")

				// Declare exchanges and queues
				log.Printf("[RABBITMQ] Setting up queues...")
				if err := setupQueues(); err == nil {
					log.Printf("[RABBITMQ SUCCESS] RabbitMQ connected and configured successfully")
					log.Printf("[RABBITMQ] Ready to publish and consume messages")
					return nil
				} else {
					log.Printf("[RABBITMQ ERROR] Failed to setup queues: %v", err)
					Conn.Close()
				}
			} else {
				log.Printf("[RABBITMQ ERROR] Failed to open channel: %v", err)
				Conn.Close()
			}
		}

		if i < maxRetries-1 {
			log.Printf("[RABBITMQ] Connection failed: %v. Retrying in %v...", err, retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Printf("[RABBITMQ ERROR] Failed to connect after %d attempts", maxRetries)
	return fmt.Errorf("failed to connect to RabbitMQ after %d attempts", maxRetries)
}

// Close closes RabbitMQ connection
func Close() error {
	log.Printf("[RABBITMQ] Closing RabbitMQ connection...")
	if Channel != nil {
		log.Printf("[RABBITMQ] Closing channel...")
		if err := Channel.Close(); err != nil {
			log.Printf("[RABBITMQ ERROR] Error closing channel: %v", err)
		} else {
			log.Printf("[RABBITMQ] Channel closed successfully")
		}
	}
	if Conn != nil {
		log.Printf("[RABBITMQ] Closing connection...")
		if err := Conn.Close(); err != nil {
			log.Printf("[RABBITMQ ERROR] Error closing connection: %v", err)
			return err
		}
		log.Printf("[RABBITMQ] Connection closed successfully")
	}
	return nil
}

// setupQueues declares all required queues
func setupQueues() error {
	// Email queue for sending emails
	log.Printf("[RABBITMQ] Declaring queue: email_queue")
	emailQueue, err := Channel.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare email_queue: %w", err)
	}
	log.Printf("[RABBITMQ] Queue declared: email_queue (Messages: %d, Consumers: %d)",
		emailQueue.Messages, emailQueue.Consumers)

	// Notification queue
	log.Printf("[RABBITMQ] Declaring queue: notification_queue")
	notifQueue, err := Channel.QueueDeclare(
		"notification_queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare notification_queue: %w", err)
	}
	log.Printf("[RABBITMQ] Queue declared: notification_queue (Messages: %d, Consumers: %d)",
		notifQueue.Messages, notifQueue.Consumers)

	// Feed processing queue
	log.Printf("[RABBITMQ] Declaring queue: feed_processing_queue")
	feedQueue, err := Channel.QueueDeclare(
		"feed_processing_queue", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare feed_processing_queue: %w", err)
	}
	log.Printf("[RABBITMQ] Queue declared: feed_processing_queue (Messages: %d, Consumers: %d)",
		feedQueue.Messages, feedQueue.Consumers)

	// User activity queue for tracking user events
	log.Printf("[RABBITMQ] Declaring queue: user_activity_queue")
	activityQueue, err := Channel.QueueDeclare(
		"user_activity_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare user_activity_queue: %w", err)
	}
	log.Printf("[RABBITMQ] Queue declared: user_activity_queue (Messages: %d, Consumers: %d)",
		activityQueue.Messages, activityQueue.Consumers)

	log.Printf("[RABBITMQ] All queues declared successfully")
	return nil
}
