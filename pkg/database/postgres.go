package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"lostmediago/internal/config"
	"lostmediago/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection with polling retry
func Connect() error {
	cfg := config.AppConfig.Database

	log.Printf("[DATABASE] Connecting to PostgreSQL...")
	log.Printf("[DATABASE] Host: %s, Port: %s, Database: %s, User: %s, SSLMode: %s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.SSLMode)

	var err error
	var dsn string

	if cfg.URL != "" {
		// Ensure sslmode is in URL if not present
		dsn = cfg.URL
		if cfg.SSLMode != "" && !strings.Contains(dsn, "sslmode=") {
			separator := "?"
			if strings.Contains(dsn, "?") {
				separator = "&"
			}
			dsn = fmt.Sprintf("%s%ssslmode=%s", dsn, separator, cfg.SSLMode)
		} else if cfg.SSLMode != "" && strings.Contains(dsn, "sslmode=") {
			// Replace existing sslmode if present
			dsn = replaceSSLMode(dsn, cfg.SSLMode)
		}
		log.Printf("[DATABASE] Using connection URL (SSLMode: %s)", cfg.SSLMode)
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)
		log.Printf("[DATABASE] Using connection string")
	}

	// Poll until database is ready
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("[DATABASE] Attempting connection (attempt %d/%d)...", i+1, maxRetries)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true, // Disable FK constraint during migration
		})

		if err == nil {
			// Test connection
			sqlDB, err := DB.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					log.Printf("[DATABASE SUCCESS] Connected to PostgreSQL successfully")

					// Set connection pool settings
					sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
					sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
					sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

					log.Printf("[DATABASE] Connection pool configured - MaxOpen: %d, MaxIdle: %d, MaxLifetime: %v",
						cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)
					return nil
				}
			}
		}

		if i < maxRetries-1 {
			log.Printf("[DATABASE] Connection failed: %v. Retrying in %v...", err, retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Printf("[DATABASE ERROR] Failed to connect after %d attempts: %v", maxRetries, err)
	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// replaceSSLMode replaces sslmode in connection string
func replaceSSLMode(dsn, sslMode string) string {
	if strings.Contains(dsn, "sslmode=") {
		parts := strings.SplitN(dsn, "sslmode=", 2)
		if len(parts) == 2 {
			afterSSL := parts[1]
			nextParam := strings.Index(afterSSL, "&")
			if nextParam != -1 {
				return parts[0] + "sslmode=" + sslMode + afterSSL[nextParam:]
			}
			return parts[0] + "sslmode=" + sslMode
		}
	}
	return dsn
}

// Close closes database connection
func Close() error {
	if DB != nil {
		log.Printf("[DATABASE] Closing database connection...")
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("[DATABASE ERROR] Error getting underlying sql.DB: %v", err)
			return err
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("[DATABASE ERROR] Error closing database connection: %v", err)
			return err
		}
		log.Printf("[DATABASE] Database connection closed successfully")
		return nil
	}
	return nil
}

// AutoMigrate runs GORM auto-migration for all models
// Migrates in order: parent tables first, then child tables
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	log.Printf("[DATABASE] Running auto-migration...")

	// Migrate in order: parent tables first, then child tables with foreign keys
	// This ensures all referenced tables exist before creating foreign key constraints

	// Step 1: Migrate parent/standalone tables first
	log.Printf("[DATABASE] Step 1: Migrating parent tables (users, roles)...")
	parentInstances := []interface{}{
		&models.User{},
		&models.Role{},
	}
	if err := DB.AutoMigrate(parentInstances...); err != nil {
		log.Printf("[DATABASE ERROR] Parent tables migration failed: %v", err)
		return fmt.Errorf("parent tables migration failed: %w", err)
	}
	log.Printf("[DATABASE] ✓ Parent tables migrated successfully")

	// Step 2: Migrate tables that reference parent tables
	log.Printf("[DATABASE] Step 2: Migrating child tables (posts, comments, likes, etc.)...")
	childInstances := []interface{}{
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Follower{},
		&models.Message{},
		&models.Notification{},
		&models.ContentSection{},
		&models.Payment{},
	}
	if err := DB.AutoMigrate(childInstances...); err != nil {
		log.Printf("[DATABASE ERROR] Child tables migration failed: %v", err)
		return fmt.Errorf("child tables migration failed: %w", err)
	}
	log.Printf("[DATABASE] ✓ Child tables migrated successfully")

	log.Printf("[DATABASE] Auto-migration completed successfully")

	// Create composite indexes that GORM doesn't support directly
	log.Printf("[DATABASE] Creating composite indexes...")
	if err := createCompositeIndexes(); err != nil {
		log.Printf("[DATABASE WARNING] Some composite indexes may not be created: %v", err)
		// Don't fail migration if indexes fail, just warn
	} else {
		log.Printf("[DATABASE] Composite indexes created successfully")
	}

	return nil
}

// createCompositeIndexes creates composite indexes that GORM doesn't support in tags
func createCompositeIndexes() error {
	indexes := []struct {
		name    string
		table   string
		columns string
		unique  bool
	}{
		// Users table (GORM uses snake_case for column names)
		{"idx_user_lookup", "users", "username, email", false},

		// Posts table
		{"idx_post_category", "posts", "category, is_published, is_deleted", false},
		{"idx_post_feed", "posts", "created_at, is_published", false},
		{"idx_scheduled_posts", "posts", "scheduled_at, is_scheduled", false},

		// Comments table
		{"idx_post_comments", "comments", "post_id, is_deleted, created_at", false},
		{"idx_user_comments", "comments", "user_id, created_at", false},
		{"idx_comment_replies", "comments", "parent_id, is_deleted", false},

		// Likes table
		{"idx_user_post_like", "likes", "user_id, post_id", true},
		{"idx_user_comment_like", "likes", "user_id, comment_id", true},
		{"idx_post_likes", "likes", "post_id, created_at", false},
		{"idx_comment_likes", "likes", "comment_id, created_at", false},

		// Followers table
		{"idx_unique_follow", "followers", "follower_id, following_id", true},
		{"idx_followers_list", "followers", "following_id, is_active", false},
		{"idx_following_list", "followers", "follower_id, is_active", false},

		// Messages table
		{"idx_inbox", "messages", "receiver_id, is_read, created_at", false},
		{"idx_conversation", "messages", "sender_id, receiver_id, created_at", false},
		{"idx_deleted_messages", "messages", "receiver_id, is_deleted", false},

		// Notifications table
		{"idx_user_notifications", "notifications", "user_id, is_read, created_at", false},
		{"idx_actor_notifications", "notifications", "actor_id, created_at", false},

		// Content sections table
		{"idx_post_sections", "content_sections", "post_id, \"order\"", false},

		// Payments table
		{"idx_user_payments", "payments", "user_id, status, created_at", false},
		{"idx_order_status", "payments", "order_id, status", false},
		{"idx_pending_payments", "payments", "status, expiry_time", false},
	}

	for _, idx := range indexes {
		uniqueClause := ""
		if idx.unique {
			uniqueClause = "UNIQUE "
		}

		// Check if index already exists
		var exists bool
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM pg_indexes 
				WHERE indexname = $1 AND tablename = $2
			)
		`
		if err := DB.Raw(checkQuery, idx.name, idx.table).Scan(&exists).Error; err != nil {
			log.Printf("[DATABASE] Error checking index %s: %v", idx.name, err)
			continue
		}

		if exists {
			log.Printf("[DATABASE] Index %s already exists, skipping", idx.name)
			continue
		}

		// Create index
		createQuery := fmt.Sprintf(
			"CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
			uniqueClause,
			idx.name,
			idx.table,
			idx.columns,
		)

		if err := DB.Exec(createQuery).Error; err != nil {
			log.Printf("[DATABASE] Error creating index %s: %v", idx.name, err)
			// Continue with other indexes
			continue
		}

		log.Printf("[DATABASE] Created index: %s on %s", idx.name, idx.table)
	}

	return nil
}
