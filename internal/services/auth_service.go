package services

import (
	"errors"
	"log"
	"lostmediago/internal/models"
	"lostmediago/internal/repositories"
	"lostmediago/internal/utils"
	"lostmediago/pkg/mq"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, string, error)
	Login(req *models.LoginRequest) (*models.User, string, string, error)
	GoogleOAuth(req *models.GoogleOAuthRequest) (*models.User, string, string, error)
	RefreshToken(refreshToken string) (*models.User, string, string, error)
	VerifyEmail(token string) (*models.User, string, string, error)
	VerifyOTP(email, otpCode string) (*models.User, string, string, error)
	ForgotPassword(email string) error
	VerifyResetPassword(token string) error
	ResetPassword(token, newPassword string) (*models.User, string, string, error)
	UpdateProfile(userId string, req *models.UpdateProfileRequest) (*models.User, error)
	GetUserByID(userId string) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, string, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		// Check if registered with Google
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.GoogleId != nil && *existingUser.GoogleId != "" {
			return nil, "", errors.New("email already registered with Google. Please use Google sign in")
		}

		// If email exists but not verified, resend verification email
		if existingUser != nil && !existingUser.IsEmailVerified {
			// Generate new OTP code (6 digits)
			otpCode, err := utils.GenerateOTP()
			if err != nil {
				return nil, "", errors.New("failed to generate OTP code")
			}

			expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes expiration for OTP

			// Update verification token with OTP code
			err = s.userRepo.UpdateEmailVerificationToken(existingUser.UserId, otpCode, expiresAt)
			if err != nil {
				return nil, "", errors.New("failed to update verification token")
			}

			// Resend verification email with OTP (async via message broker)
			err = mq.PublishVerificationEmail(existingUser.Email, otpCode)
			if err != nil {
				// Log error but don't block registration
				// Email will be retried by worker if RabbitMQ recovers
				log.Printf("[WARNING] Failed to publish verification email to queue: %v", err)
			}

			// Return existing user with new OTP code (so frontend can redirect to OTP page)
			return existingUser, otpCode, nil
		}

		// Email already registered and verified
		return nil, "", errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// Generate OTP code (6 digits)
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return nil, "", errors.New("failed to generate OTP code")
	}

	expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes expiration for OTP

	// Create user
	user := &models.User{
		UserId:                   uuid.New().String(),
		Username:                 req.Username,
		Email:                    req.Email,
		Password:                 &hashedPassword,
		Role:                     "member",
		Star:                     0,
		IsEmailVerified:          false,
		EmailVerificationToken:   &otpCode,
		EmailVerificationExpires: &expiresAt,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}

	// Publish verification email to queue (async, non-blocking)
	err = mq.PublishVerificationEmail(user.Email, otpCode)
	if err != nil {
		// Log error but don't block registration
		// Email will be retried by worker if RabbitMQ recovers
		log.Printf("[WARNING] Failed to publish verification email to queue: %v", err)
	}

	// Publish registration event (async, non-blocking)
	_ = mq.PublishRegisterEvent(user.UserId, user.Email)

	return user, otpCode, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.User, string, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, "", "", errors.New("account is banned")
	}

	// Check if user is registered with Google (no password)
	if user.Password == nil || *user.Password == "" {
		return nil, "", "", errors.New("email registered with Google. Please use Google sign in")
	}

	if !utils.ComparePassword(req.Password, *user.Password) {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Check if email is not verified
	if !user.IsEmailVerified {
		// Generate new OTP code (6 digits)
		otpCode, err := utils.GenerateOTP()
		if err != nil {
			return nil, "", "", errors.New("failed to generate OTP code")
		}

		expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes expiration for OTP

		// Update verification token with OTP code
		err = s.userRepo.UpdateEmailVerificationToken(user.UserId, otpCode, expiresAt)
		if err != nil {
			return nil, "", "", errors.New("failed to update verification token")
		}

		// Send verification email with OTP (async via message broker)
		err = mq.PublishVerificationEmail(user.Email, otpCode)
		if err != nil {
			// Log error but don't block login
			// Email will be retried by worker if RabbitMQ recovers
			log.Printf("[WARNING] Failed to publish verification email to queue: %v", err)
		}

		// Return error indicating email needs verification (special error code)
		return user, otpCode, "", errors.New("EMAIL_NOT_VERIFIED")
	}

	// Generate tokens
	token, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Publish login event (async via message broker) - includes UpdateLastLogin
	_ = mq.PublishLoginEvent(user.UserId, user.Email)

	return user, token, refreshToken, nil
}

func (s *authService) GoogleOAuth(req *models.GoogleOAuthRequest) (*models.User, string, string, error) {
	// Check if user exists by Google ID first (faster lookup)
	existingUser, err := s.userRepo.FindByGoogleID(req.GoogleId)
	if err == nil && existingUser != nil {
		// User exists with this Google ID - login
		// Publish login event (async via message broker) - includes UpdateLastLogin
		_ = mq.PublishLoginEvent(existingUser.UserId, existingUser.Email)

		// Generate tokens
		token, err := utils.GenerateToken(existingUser.UserId, existingUser.Email, existingUser.Role)
		if err != nil {
			return nil, "", "", errors.New("failed to generate token")
		}

		refreshToken, err := utils.GenerateRefreshToken(existingUser.UserId, existingUser.Email, existingUser.Role)
		if err != nil {
			return nil, "", "", errors.New("failed to generate refresh token")
		}

		return existingUser, token, refreshToken, nil
	}

	// Check if user exists by email
	existingUser, err = s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		// User exists - check if already registered with password
		if existingUser.Password != nil && *existingUser.Password != "" {
			return nil, "", "", errors.New("email already registered with password. Please use email and password to login")
		}

		// User exists but no password - check Google ID
		if existingUser.GoogleId != nil && *existingUser.GoogleId != req.GoogleId {
			return nil, "", "", errors.New("email already registered with different Google account")
		}

		// Update user with Google ID if not set
		if existingUser.GoogleId == nil {
			existingUser.GoogleId = &req.GoogleId
			if req.ProfilePhoto != "" {
				existingUser.ProfilePic = &req.ProfilePhoto
			}
			existingUser.IsEmailVerified = true // Google emails are verified
			existingUser.UpdatedAt = time.Now()
			if err := s.userRepo.Update(existingUser); err != nil {
				return nil, "", "", errors.New("failed to update user")
			}
		}

		// Publish login event (async via message broker) - includes UpdateLastLogin
		_ = mq.PublishLoginEvent(existingUser.UserId, existingUser.Email)

		// Generate tokens
		token, err := utils.GenerateToken(existingUser.UserId, existingUser.Email, existingUser.Role)
		if err != nil {
			return nil, "", "", errors.New("failed to generate token")
		}

		refreshToken, err := utils.GenerateRefreshToken(existingUser.UserId, existingUser.Email, existingUser.Role)
		if err != nil {
			return nil, "", "", errors.New("failed to generate refresh token")
		}

		return existingUser, token, refreshToken, nil
	}

	// New user - create account with Google
	// Generate username from email
	username := req.Email
	if atIndex := strings.Index(username, "@"); atIndex != -1 {
		username = username[:atIndex]
	}

	// Create user
	user := &models.User{
		UserId:          uuid.New().String(),
		GoogleId:        &req.GoogleId,
		Username:        username,
		Email:           req.Email,
		Password:        nil, // No password for Google OAuth users
		ProfilePic:      &req.ProfilePhoto,
		Role:            "member",
		Star:            0,
		IsEmailVerified: true, // Google emails are pre-verified
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", "", errors.New("failed to create user")
	}

	// Publish registration event for Google OAuth (async via message broker)
	_ = mq.PublishRegisterEvent(user.UserId, user.Email)

	// Generate tokens
	token, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Publish login event (async via message broker) - includes UpdateLastLogin
	_ = mq.PublishLoginEvent(user.UserId, user.Email)

	return user, token, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (*models.User, string, string, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired refresh token")
	}

	// Find user
	user, err := s.userRepo.FindByID(claims.UserId)
	if err != nil {
		return nil, "", "", errors.New("user not found")
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, "", "", errors.New("account is banned")
	}

	// Generate new tokens
	token, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	return user, token, newRefreshToken, nil
}

func (s *authService) VerifyEmail(token string) (*models.User, string, string, error) {
	// Find user by verification token
	user, err := s.userRepo.FindByEmailVerificationToken(token)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired verification token")
	}

	// Update email verification status
	err = s.userRepo.UpdateEmailVerification(user.UserId, true)
	if err != nil {
		return nil, "", "", errors.New("failed to verify email")
	}

	// Publish email verified event (async via message broker)
	_ = mq.PublishEmailVerifiedEvent(user.UserId, user.Email)

	// Generate JWT tokens for auto login
	accessToken, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Publish login event (async via message broker) - includes UpdateLastLogin
	_ = mq.PublishLoginEvent(user.UserId, user.Email)

	return user, accessToken, refreshToken, nil
}

func (s *authService) VerifyOTP(email, otpCode string) (*models.User, string, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", "", errors.New("user not found")
	}

	// Check if OTP code matches
	if user.EmailVerificationToken == nil || *user.EmailVerificationToken != otpCode {
		return nil, "", "", errors.New("invalid OTP code")
	}

	// Check if OTP is expired
	if user.EmailVerificationExpires == nil || time.Now().After(*user.EmailVerificationExpires) {
		return nil, "", "", errors.New("OTP code has expired")
	}

	// Check if already verified
	if user.IsEmailVerified {
		return nil, "", "", errors.New("email already verified")
	}

	// Update email verification status
	err = s.userRepo.UpdateEmailVerification(user.UserId, true)
	if err != nil {
		return nil, "", "", errors.New("failed to verify email")
	}

	// Publish email verified event (async via message broker)
	_ = mq.PublishEmailVerifiedEvent(user.UserId, user.Email)

	// Generate JWT tokens for auto login
	accessToken, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Publish login event (async via message broker) - includes UpdateLastLogin
	_ = mq.PublishLoginEvent(user.UserId, user.Email)

	return user, accessToken, refreshToken, nil
}

func (s *authService) ForgotPassword(email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if email exists for security
		return nil
	}

	// Generate password reset token
	resetToken, err := utils.GeneratePasswordResetToken()
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	expiresAt := time.Now().Add(1 * time.Hour) // 1 hour expiration

	// Update user with reset token
	err = s.userRepo.UpdatePasswordResetToken(user.UserId, resetToken, expiresAt)
	if err != nil {
		return errors.New("failed to set reset token")
	}

	// Publish password reset email to queue (async, non-blocking)
	err = mq.PublishPasswordResetEmail(user.Email, resetToken)
	if err != nil {
		// Log error but don't block
		// Email will be retried by worker if RabbitMQ recovers
		log.Printf("[WARNING] Failed to publish password reset email to queue: %v", err)
	}

	return nil
}

func (s *authService) VerifyResetPassword(token string) error {
	// Find user by reset token
	_, err := s.userRepo.FindByPasswordResetToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	return nil
}

func (s *authService) ResetPassword(token, newPassword string) (*models.User, string, string, error) {
	// Find user by reset token
	user, err := s.userRepo.FindByPasswordResetToken(token)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, "", "", errors.New("failed to hash password")
	}

	// Update password and clear reset token
	user.Password = &hashedPassword
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, "", "", errors.New("failed to update password")
	}

	err = s.userRepo.ClearPasswordResetToken(user.UserId)
	if err != nil {
		return nil, "", "", errors.New("failed to clear reset token")
	}

	// Publish password reset event (async via message broker)
	_ = mq.PublishPasswordResetEvent(user.UserId, user.Email)

	// Generate JWT tokens for auto login
	accessToken, err := utils.GenerateToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Publish login event (async via message broker) - includes UpdateLastLogin
	_ = mq.PublishLoginEvent(user.UserId, user.Email)

	return user, accessToken, refreshToken, nil
}

func (s *authService) UpdateProfile(userId string, req *models.UpdateProfileRequest) (*models.User, error) {
	// Get user by ID
	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update username if provided (allow duplicate usernames across users)
	if req.Username != nil && *req.Username != user.Username {
		user.Username = *req.Username
	}

	// Update bio if provided
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	// Update profile picture if provided
	if req.ProfilePic != nil {
		user.ProfilePic = req.ProfilePic
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}

func (s *authService) GetUserByID(userId string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
