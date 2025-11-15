package utils

import (
	"lostmediago/internal/models"
)

// ConvertUserToResponse converts User model to UserResponse
func ConvertUserToResponse(user *models.User) *models.UserResponse {
	loginType := "credential"
	if user.GoogleId != nil && *user.GoogleId != "" {
		loginType = "google"
	}

	profilePhoto := ""
	if user.ProfilePic != nil {
		profilePhoto = *user.ProfilePic
	}

	// Use username as fullName (can be enhanced with separate full_name field later)
	fullName := user.Username

	return &models.UserResponse{
		ID:           user.UserId,
		Username:     user.Username,
		Email:        user.Email,
		FullName:     fullName,
		ProfilePhoto: profilePhoto,
		IsVerified:   user.IsEmailVerified,
		UserType:     user.Role,
		LoginType:    loginType,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ConvertToAuthResponse converts User and tokens to AuthResponse
func ConvertToAuthResponse(user *models.User, accessToken, refreshToken string, expiresIn int) *models.AuthResponse {
	return &models.AuthResponse{
		User:                 ConvertUserToResponse(user),
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		ExpiresIn:            expiresIn,
		RequiresVerification: false, // Default to false, will be set explicitly when needed
	}
}
