package handlers

import (
	"net/http"

	"lostmediago/internal/models"
	"lostmediago/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecases.AuthUsecase
}

func NewAuthHandler(authUsecase usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.Register(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		if errorMessage == "email already registered" {
			statusCode = http.StatusConflict
			errorCode = "CONFLICT"
		}

		if errorMessage == "email already registered with Google. Please use Google sign in" {
			statusCode = http.StatusConflict
			errorCode = "CONFLICT"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Data: response,
	})
}

// Login handles user login
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.Login(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		if errorMessage == "invalid email or password" || errorMessage == "account is banned" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		if errorMessage == "email registered with Google. Please use Google sign in" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	// If requires verification, return 200 with verification info (not an error)
	if response.RequiresVerification {
		c.JSON(http.StatusOK, Response{
			Data: response,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// GoogleOAuth handles Google OAuth login/register
// @Summary Google OAuth login/register
// @Description Login or register using Google OAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.GoogleOAuthRequest true "Google OAuth request"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/auth/google-oauth [post]
func (h *AuthHandler) GoogleOAuth(c *gin.Context) {
	var req models.GoogleOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.GoogleOAuth(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		if errorMessage == "email already registered with password. Please use email and password to login" {
			statusCode = http.StatusConflict
			errorCode = "CONFLICT"
		}

		if errorMessage == "email already registered with different Google account" {
			statusCode = http.StatusConflict
			errorCode = "CONFLICT"
		}

		if errorMessage == "email registered with Google. Please use Google sign in" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.RefreshToken(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		if errorMessage == "invalid or expired refresh token" || errorMessage == "account is banned" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// VerifyEmail handles email verification
// @Summary Verify email address
// @Description Verify user email address with verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyEmailRequest true "Verify email request"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.VerifyEmail(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			errorMessage = validationErr.Message
		}

		if errorMessage == "invalid or expired verification token" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// VerifyOTP handles OTP verification
// @Summary Verify OTP
// @Description Verify user email with OTP code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Verify OTP request"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.VerifyOTP(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		if errorMessage == "invalid OTP code" || errorMessage == "OTP code has expired" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		if errorMessage == "email already verified" {
			statusCode = http.StatusBadRequest
			errorCode = "ALREADY_VERIFIED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// ForgotPassword handles forgot password request
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	err := h.authUsecase.ForgotPassword(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
			errorMessage = validationErr.Message
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	// Always return success for security (don't reveal if email exists)
	c.JSON(http.StatusOK, Response{
		Data: gin.H{
			"message": "If the email exists, a password reset link has been sent",
		},
	})
}

// VerifyResetPassword verifies the password reset token
// @Summary Verify password reset token
// @Description Verify if password reset token is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyResetPasswordRequest true "Verify reset password request"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/verify-reset-password [post]
func (h *AuthHandler) VerifyResetPassword(c *gin.Context) {
	var req models.VerifyResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	err := h.authUsecase.VerifyResetPassword(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			errorMessage = validationErr.Message
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: gin.H{
			"message": "Reset token is valid",
		},
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset user password with reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.authUsecase.ResetPassword(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		errorMessage := err.Error()

		if validationErr, ok := err.(*usecases.ValidationError); ok {
			errorMessage = validationErr.Message
		}

		if errorMessage == "invalid or expired reset token" {
			statusCode = http.StatusUnauthorized
			errorCode = "UNAUTHORIZED"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// Response represents a success response
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
