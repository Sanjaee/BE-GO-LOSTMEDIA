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
	"lostmediago/pkg/database"
	"lostmediago/pkg/mq"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize repositories (needed before workers)
	userRepo := repositories.NewUserRepository()

	// Connect to RabbitMQ
	if err := mq.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v. Events will not be processed.", err)
	} else {
		defer mq.Close()
		// Start email worker
		if err := workers.StartEmailWorker(); err != nil {
			log.Printf("Warning: Failed to start email worker: %v. Email events will not be processed.", err)
		}
		// Start user activity worker (handles UpdateLastLogin, analytics, etc.)
		if err := workers.StartUserActivityWorker(userRepo); err != nil {
			log.Printf("Warning: Failed to start user activity worker: %v. User activity events will not be processed.", err)
		}
	}

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize use cases
	authUsecase := usecases.NewAuthUsecase(authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUsecase)

	// Setup routes
	router := delivery.SetupRoutes(authHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Server.Port),
		Handler:        router,
		ReadTimeout:    config.AppConfig.Server.ReadTimeout,
		WriteTimeout:   config.AppConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
