package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId                   string     `json:"userId" gorm:"primaryKey;type:varchar(36)"`
	GoogleId                 *string    `json:"googleId,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	Username                 string     `json:"username" gorm:"type:varchar(255);not null;index:idx_user_lookup"`
	Email                    string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null;index:idx_user_lookup"`
	Password                 *string    `json:"-" gorm:"type:varchar(255)"`
	ProfilePic               *string    `json:"profilePic,omitempty" gorm:"type:text"`
	Bio                      *string    `json:"bio,omitempty" gorm:"type:text"`
	CreatedAt                time.Time  `json:"createdAt" gorm:"default:now()"`
	UpdatedAt                time.Time  `json:"updatedAt" gorm:"default:now()"`
	FollowersCount           int        `json:"followersCount" gorm:"default:0"`
	FollowingCount           int        `json:"followingCount" gorm:"default:0"`
	Role                     string     `json:"role" gorm:"type:varchar(50);default:'member'"`
	Star                     int        `json:"star" gorm:"default:0"`
	IsBanned                 bool       `json:"isBanned" gorm:"default:false;index:idx_banned_users"`
	BanReason                *string    `json:"banReason,omitempty" gorm:"type:text"`
	BannedBy                 *string    `json:"bannedBy,omitempty" gorm:"type:varchar(36)"`
	PostsCount               int        `json:"postsCount" gorm:"default:0"`
	IsEmailVerified          bool       `json:"isEmailVerified" gorm:"default:false;index:idx_email_verified"`
	EmailVerificationToken   *string    `json:"-" gorm:"type:varchar(255);index:idx_email_verify_token"`
	EmailVerificationExpires *time.Time `json:"-" gorm:"type:timestamp"`
	PasswordResetToken       *string    `json:"-" gorm:"type:varchar(255);index:idx_password_reset_token"`
	PasswordResetExpires     *time.Time `json:"-" gorm:"type:timestamp"`
	LastLoginAt              *time.Time `json:"lastLoginAt,omitempty" gorm:"type:timestamp"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserId == "" {
		u.UserId = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook to update UpdatedAt
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleOAuthRequest struct {
	Email        string `json:"email" binding:"required,email"`
	FullName     string `json:"full_name" binding:"required"`
	ProfilePhoto string `json:"profile_photo"`
	GoogleId     string `json:"google_id" binding:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required,len=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetPasswordRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

type AuthResponse struct {
	User                 *UserResponse `json:"user"`
	AccessToken          string        `json:"access_token"`
	RefreshToken         string        `json:"refresh_token"`
	ExpiresIn            int           `json:"expires_in"`
	RequiresVerification bool          `json:"requires_verification,omitempty"`
	VerificationToken    string        `json:"verification_token,omitempty"`
}

type UserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name,omitempty"`
	ProfilePhoto string `json:"profile_photo,omitempty"`
	IsVerified   bool   `json:"is_verified"`
	UserType     string `json:"user_type,omitempty"`
	LoginType    string `json:"login_type"`
	CreatedAt    string `json:"created_at"`
}

type RegisterResponse struct {
	User                 *UserResponse `json:"user"`
	VerificationToken    string        `json:"verification_token,omitempty"`
	RequiresVerification bool          `json:"requires_verification"`
	Message              string        `json:"message"`
}

type TokenPayload struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
