package middleware

import (
	"lostmediago/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets up CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     []string{config.AppConfig.Server.ClientURL, config.AppConfig.Server.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(cfg)
}
