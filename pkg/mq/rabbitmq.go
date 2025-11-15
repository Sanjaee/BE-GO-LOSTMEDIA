package mq

import (
	"fmt"
	"log"
	"lostmediago/internal/config"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// Connect initializes RabbitMQ connection
func Connect() error {
	cfg := config.AppConfig

	var dsn string
	if cfg.RabbitMQ.VHost != "" && cfg.RabbitMQ.VHost != "/" {
		dsn = fmt.Sprintf("amqp://%s:%s@%s:%d%s",
			cfg.RabbitMQ.User,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Host,
			cfg.RabbitMQ.Port,
			cfg.RabbitMQ.VHost,
		)
	} else {
		dsn = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			cfg.RabbitMQ.User,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Host,
			cfg.RabbitMQ.Port,
		)
	}

	var err error
	Conn, err = amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchanges and queues
	if err := setupQueues(); err != nil {
		return fmt.Errorf("failed to setup queues: %w", err)
	}

	log.Println("RabbitMQ connected successfully")
	return nil
}

// Close closes RabbitMQ connection
func Close() error {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		return Conn.Close()
	}
	return nil
}

// setupQueues declares all required queues
func setupQueues() error {
	// Email queue for sending emails
	_, err := Channel.QueueDeclare(
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

	// Notification queue
	_, err = Channel.QueueDeclare(
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

	// Feed processing queue
	_, err = Channel.QueueDeclare(
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

	// User activity queue for tracking user events
	_, err = Channel.QueueDeclare(
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

	return nil
}
