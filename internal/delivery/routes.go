package delivery

import (
	"lostmediago/internal/handlers"
	"lostmediago/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *handlers.AuthHandler, postHandler *handlers.PostHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "lostmediago-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google-oauth", authHandler.GoogleOAuth)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/verify-reset-password", authHandler.VerifyResetPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Post routes
		posts := v1.Group("/posts")
		{
			// Public routes
			posts.GET("", middleware.OptionalAuthMiddleware(), postHandler.GetAllPosts)
			posts.GET("/:id", middleware.OptionalAuthMiddleware(), postHandler.GetPost)

			// Protected routes (require authentication)
			postsProtected := posts.Group("", middleware.AuthMiddleware())
			{
				postsProtected.POST("", postHandler.CreatePost)
				postsProtected.PUT("/:id", postHandler.UpdatePost)
				postsProtected.DELETE("/:id", postHandler.DeletePost)
				postsProtected.POST("/:id/like", postHandler.LikePost)
				postsProtected.GET("/user/posts-count", postHandler.GetUserPostsCount)
			}
		}
	}

	return router
}
