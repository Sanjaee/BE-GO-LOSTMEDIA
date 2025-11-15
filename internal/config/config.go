package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	JWT      JWTConfig
	Google   GoogleConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	Env          string
	ClientURL    string
	FrontendURL  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	URL             string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type EmailConfig struct {
	From     string
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

var AppConfig *Config

func Load() error {
	AppConfig = &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "5000"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			Env:          getEnv("NODE_ENV", "development"),
			ClientURL:    getEnv("CLIENT_URL", "http://localhost:3000"),
			FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnv("POSTGRES_PORT", "5432"),
			User:            getEnv("POSTGRES_USER", "lostmedia_db"),
			Password:        getEnv("POSTGRES_PASSWORD", "123321"),
			DBName:          getEnv("POSTGRES_DB", "lostmedia"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: 5 * time.Minute,
			URL:             getEnv("DATABASE_URL", "postgresql://lostmedia_db:123321@db:5432/lostmedia"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
			User:     getEnv("RABBITMQ_USER", "lostmediago"),
			Password: getEnv("RABBITMQ_PASSWORD", "password123"),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "D8D3DA7A75F61ACD5A4CD579EDBBC"),
			Expiration:        24 * time.Hour,
			RefreshExpiration: 168 * time.Hour, // 7 days
		},
		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:5000/api/v1/auth/google/callback"),
		},
		Email: EmailConfig{
			From:     getEnv("EMAIL_FROM", "noreply@lostmediago.com"),
			SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getEnvAsInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

