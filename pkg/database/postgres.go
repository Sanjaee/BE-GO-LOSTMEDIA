package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lostmediago/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes database connection
func Connect() error {
	cfg := config.AppConfig.Database

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
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)
	}

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
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
		return DB.Close()
	}
	return nil
}

// WithTransaction executes a function within a database transaction
func WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// NullTime converts time.Time to sql.NullTime
func NullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// NullString converts string to sql.NullString
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	if *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
