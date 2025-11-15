package usecases

import (
	"lostmediago/internal/models"
	"lostmediago/internal/services"
	"lostmediago/internal/utils"
)

type AuthUsecase interface {
	Register(req *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	GoogleOAuth(req *models.GoogleOAuthRequest) (*models.AuthResponse, error)
	RefreshToken(req *models.RefreshTokenRequest) (*models.AuthResponse, error)
	VerifyEmail(req *models.VerifyEmailRequest) (*models.AuthResponse, error)
	VerifyOTP(req *models.VerifyOTPRequest) (*models.AuthResponse, error)
	ForgotPassword(req *models.ForgotPasswordRequest) error
	VerifyResetPassword(req *models.VerifyResetPasswordRequest) error
	ResetPassword(req *models.ResetPasswordRequest) (*models.AuthResponse, error)
}

type authUsecase struct {
	authService services.AuthService
}

func NewAuthUsecase(authService services.AuthService) AuthUsecase {
	return &authUsecase{
		authService: authService,
	}
}

func (uc *authUsecase) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, &ValidationError{Message: "all fields are required"}
	}

	// Register user (or resend OTP if email exists but not verified)
	user, verificationToken, err := uc.authService.Register(req)
	if err != nil {
		return nil, err
	}

	// Check if user requires verification
	requiresVerification := !user.IsEmailVerified

	response := &models.RegisterResponse{
		User:                 utils.ConvertUserToResponse(user),
		VerificationToken:    verificationToken, // Only included if requires verification
		RequiresVerification: requiresVerification,
		Message:              "Registration successful. Please verify your email.",
	}

	return response, nil
}

func (uc *authUsecase) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, &ValidationError{Message: "email and password are required"}
	}

	// Login user
	user, token, refreshToken, err := uc.authService.Login(req)
	if err != nil {
		// Check if it's an email not verified error
		if err.Error() == "EMAIL_NOT_VERIFIED" {
			// Return response with verification required
			response := &models.AuthResponse{
				User:                 utils.ConvertUserToResponse(user),
				AccessToken:          "", // No token until verified
				RefreshToken:         "",
				ExpiresIn:            0,
				RequiresVerification: true,
				VerificationToken:    token, // token here is actually verification token
			}
			return response, nil
		}
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, token, refreshToken, 86400)
	response.RequiresVerification = false

	return response, nil
}

func (uc *authUsecase) GoogleOAuth(req *models.GoogleOAuthRequest) (*models.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.GoogleId == "" || req.FullName == "" {
		return nil, &ValidationError{Message: "email, google_id, and full_name are required"}
	}

	// Google OAuth login/register
	user, token, refreshToken, err := uc.authService.GoogleOAuth(req)
	if err != nil {
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, token, refreshToken, 86400)

	return response, nil
}

func (uc *authUsecase) RefreshToken(req *models.RefreshTokenRequest) (*models.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, &ValidationError{Message: "refresh_token is required"}
	}

	// Refresh token
	user, token, refreshToken, err := uc.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, token, refreshToken, 86400)

	return response, nil
}

func (uc *authUsecase) VerifyEmail(req *models.VerifyEmailRequest) (*models.AuthResponse, error) {
	if req.Token == "" {
		return nil, &ValidationError{Message: "token is required"}
	}

	// Verify email and get JWT tokens for auto login
	user, accessToken, refreshToken, err := uc.authService.VerifyEmail(req.Token)
	if err != nil {
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, accessToken, refreshToken, 86400)

	return response, nil
}

func (uc *authUsecase) VerifyOTP(req *models.VerifyOTPRequest) (*models.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.OTPCode == "" {
		return nil, &ValidationError{Message: "email and otp_code are required"}
	}

	if len(req.OTPCode) != 6 {
		return nil, &ValidationError{Message: "otp_code must be 6 digits"}
	}

	// Verify OTP
	user, accessToken, refreshToken, err := uc.authService.VerifyOTP(req.Email, req.OTPCode)
	if err != nil {
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, accessToken, refreshToken, 86400)

	return response, nil
}

func (uc *authUsecase) ForgotPassword(req *models.ForgotPasswordRequest) error {
	if req.Email == "" {
		return &ValidationError{Message: "email is required"}
	}

	err := uc.authService.ForgotPassword(req.Email)
	return err
}

func (uc *authUsecase) VerifyResetPassword(req *models.VerifyResetPasswordRequest) error {
	if req.Token == "" {
		return &ValidationError{Message: "token is required"}
	}

	err := uc.authService.VerifyResetPassword(req.Token)
	return err
}

func (uc *authUsecase) ResetPassword(req *models.ResetPasswordRequest) (*models.AuthResponse, error) {
	if req.Token == "" || req.NewPassword == "" {
		return nil, &ValidationError{Message: "token and new password are required"}
	}

	if len(req.NewPassword) < 8 {
		return nil, &ValidationError{Message: "password must be at least 8 characters"}
	}

	// Reset password and get JWT tokens for auto login
	user, accessToken, refreshToken, err := uc.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		return nil, err
	}

	// Convert to response format (expires in 24 hours = 86400 seconds)
	response := utils.ConvertToAuthResponse(user, accessToken, refreshToken, 86400)

	return response, nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
