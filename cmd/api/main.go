package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lostmediago/internal/config"
	"lostmediago/internal/delivery"
	"lostmediago/internal/handlers"
	"lostmediago/internal/repositories"
	"lostmediago/internal/services"
	"lostmediago/internal/usecases"
	"lostmediago/internal/workers"
	"lostmediago/pkg/cache"
	"lostmediago/pkg/database"
	"lostmediago/pkg/mq"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("=========================================")
	log.Println("LostMediaGo API Server Starting...")
	log.Println("=========================================")

	// Load configuration
	log.Println("[INIT] Loading configuration...")
	if err := config.Load(); err != nil {
		log.Fatal("[INIT ERROR] Failed to load configuration:", err)
	}
	log.Printf("[INIT SUCCESS] Configuration loaded - Environment: %s", config.AppConfig.Server.Env)

	// Set Gin mode
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("[INIT] Gin mode set to: production")
	} else {
		log.Println("[INIT] Gin mode set to: development")
	}

	log.Println("-----------------------------------------")
	log.Println("[SERVICES] Initializing services...")
	log.Println("-----------------------------------------")

	// Connect to PostgreSQL database
	log.Println("[SERVICES] Connecting to PostgreSQL database...")
	if err := database.Connect(); err != nil {
		log.Fatal("[SERVICES ERROR] Failed to connect to database:", err)
	}
	defer func() {
		log.Println("[SHUTDOWN] Closing database connection...")
		database.Close()
	}()
	log.Println("[SERVICES] ✓ PostgreSQL database connected")

	// Run GORM auto-migration for all models
	// Migration order is handled inside AutoMigrate function
	log.Println("[SERVICES] Running database auto-migration...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("[SERVICES ERROR] Failed to run auto-migration:", err)
	}
	log.Println("[SERVICES] ✓ Database migration completed - All tables created")

	// Connect to Redis
	log.Println("[SERVICES] Connecting to Redis cache...")
	redisConnected := false
	if err := cache.Connect(); err != nil {
		log.Printf("[SERVICES WARNING] Failed to connect to Redis: %v", err)
		log.Printf("[SERVICES WARNING] Caching will be disabled.")
		log.Println("[SERVICES] ✗ Redis not connected")
	} else {
		redisConnected = true
		log.Println("[SERVICES] ✓ Redis connected")
		defer func() {
			log.Println("[SHUTDOWN] Closing Redis connection...")
			cache.Close()
		}()
	}

	// Initialize repositories (needed for workers and services)
	log.Println("[INIT] Initializing repositories...")
	userRepo := repositories.NewUserRepository()
	log.Println("[INIT] ✓ Repositories initialized")

	// Connect to RabbitMQ
	log.Println("[SERVICES] Connecting to RabbitMQ message broker...")
	rabbitmqConnected := false
	if err := mq.Connect(); err != nil {
		log.Printf("[SERVICES WARNING] Failed to connect to RabbitMQ: %v", err)
		log.Printf("[SERVICES WARNING] Events will not be processed. Email sending will be disabled.")
		log.Println("[SERVICES] ✗ RabbitMQ not connected")
	} else {
		rabbitmqConnected = true
		log.Println("[SERVICES] ✓ RabbitMQ connected")
		defer func() {
			log.Println("[SHUTDOWN] Closing RabbitMQ connection...")
			mq.Close()
		}()

		// Start email worker
		log.Println("[SERVICES] Starting email worker...")
		if err := workers.StartEmailWorker(); err != nil {
			log.Printf("[SERVICES WARNING] Failed to start email worker: %v", err)
			log.Printf("[SERVICES WARNING] Email events will not be processed.")
			log.Println("[SERVICES] ✗ Email worker not started")
		} else {
			log.Println("[SERVICES] ✓ Email worker started")
		}

		// Start user activity worker (handles UpdateLastLogin, analytics, etc.)
		log.Println("[SERVICES] Starting user activity worker...")
		if err := workers.StartUserActivityWorker(userRepo); err != nil {
			log.Printf("[SERVICES WARNING] Failed to start user activity worker: %v", err)
			log.Printf("[SERVICES WARNING] User activity events will not be processed.")
			log.Println("[SERVICES] ✗ User activity worker not started")
		} else {
			log.Println("[SERVICES] ✓ User activity worker started")
		}
	}

	log.Println("-----------------------------------------")
	log.Println("[SERVICES] Service Status Summary:")
	log.Println("-----------------------------------------")
	log.Println("  ✓ PostgreSQL Database")
	if redisConnected {
		log.Println("  ✓ Redis Cache")
	} else {
		log.Println("  ✗ Redis Cache (Not Connected)")
	}
	if rabbitmqConnected {
		log.Println("  ✓ RabbitMQ Message Broker")
	} else {
		log.Println("  ✗ RabbitMQ Message Broker (Not Connected)")
	}
	log.Println("-----------------------------------------")

	// Initialize services
	log.Println("[INIT] Initializing application services...")
	authService := services.NewAuthService(userRepo)
	log.Println("[INIT] ✓ Auth service initialized")

	// Initialize post repositories
	postRepo := repositories.NewPostRepository()
	likeRepo := repositories.NewLikeRepository()
	log.Println("[INIT] ✓ Post repositories initialized")

	// Initialize post service
	postService := services.NewPostService(postRepo, userRepo, likeRepo)
	log.Println("[INIT] ✓ Post service initialized")

	// Initialize use cases
	authUsecase := usecases.NewAuthUsecase(authService)
	postUsecase := usecases.NewPostUsecase(postService, userRepo)
	log.Println("[INIT] ✓ Usecases initialized")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUsecase)
	postHandler := handlers.NewPostHandler(postUsecase)
	log.Println("[INIT] ✓ Handlers initialized")

	// Setup routes
	log.Println("[INIT] Setting up routes...")
	router := delivery.SetupRoutes(authHandler, postHandler)
	log.Println("[INIT] ✓ Routes configured")

	// Create HTTP server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Server.Port),
		Handler:        router,
		ReadTimeout:    config.AppConfig.Server.ReadTimeout,
		WriteTimeout:   config.AppConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Println("=========================================")
	log.Println("[SERVER] Starting HTTP server...")
	log.Printf("[SERVER] Port: %s", config.AppConfig.Server.Port)
	log.Printf("[SERVER] Host: %s", config.AppConfig.Server.Host)
	log.Printf("[SERVER] Environment: %s", config.AppConfig.Server.Env)
	log.Printf("[SERVER] Frontend URL: %s", config.AppConfig.Server.FrontendURL)
	log.Println("=========================================")

	// Start server in goroutine
	go func() {
		log.Printf("[SERVER] HTTP server listening on :%s", config.AppConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("[SERVER ERROR] Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("=========================================")
	log.Println("[SHUTDOWN] Shutting down server...")
	log.Println("=========================================")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("[SHUTDOWN ERROR] Server forced to shutdown:", err)
	}

	log.Println("[SHUTDOWN] Server exited gracefully")
	log.Println("=========================================")
}
