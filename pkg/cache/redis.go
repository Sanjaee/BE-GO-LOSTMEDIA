package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"lostmediago/internal/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var ctx = context.Background()

// Connect initializes Redis connection with polling retry
func Connect() error {
	cfg := config.AppConfig.Redis

	log.Printf("[REDIS] Connecting to Redis...")
	log.Printf("[REDIS] Host: %s, Port: %s, DB: %d", cfg.Host, cfg.Port, cfg.DB)

	// Build Redis address
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Create Redis client
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Poll until Redis is ready
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("[REDIS] Attempting connection (attempt %d/%d)...", i+1, maxRetries)

		// Test connection
		err := Client.Ping(ctx).Err()
		if err == nil {
			log.Printf("[REDIS SUCCESS] Connected to Redis successfully")
			log.Printf("[REDIS] Connection pool configured - PoolSize: %d", cfg.PoolSize)
			return nil
		}

		if i < maxRetries-1 {
			log.Printf("[REDIS] Connection failed: %v. Retrying in %v...", err, retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Printf("[REDIS ERROR] Failed to connect after %d attempts", maxRetries)
	return fmt.Errorf("failed to connect to Redis after %d attempts", maxRetries)
}

// Close closes Redis connection
func Close() error {
	if Client != nil {
		log.Printf("[REDIS] Closing Redis connection...")
		err := Client.Close()
		if err != nil {
			log.Printf("[REDIS ERROR] Error closing Redis connection: %v", err)
			return err
		}
		log.Printf("[REDIS] Redis connection closed successfully")
		return nil
	}
	return nil
}

// Ping tests Redis connection
func Ping() error {
	if Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Client.Ping(ctx).Err()
}
