package models

import "time"

type User struct {
	UserId                   string     `json:"userId" db:"userId"`
	GoogleId                 *string    `json:"googleId,omitempty" db:"googleId"`
	Username                 string     `json:"username" db:"username"`
	Email                    string     `json:"email" db:"email"`
	Password                 *string    `json:"-" db:"password"`
	ProfilePic               *string    `json:"profilePic,omitempty" db:"profilePic"`
	Bio                      *string    `json:"bio,omitempty" db:"bio"`
	CreatedAt                time.Time  `json:"createdAt" db:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt" db:"updatedAt"`
	FollowersCount           int        `json:"followersCount" db:"followersCount"`
	FollowingCount           int        `json:"followingCount" db:"followingCount"`
	Role                     string     `json:"role" db:"role"`
	Star                     int        `json:"star" db:"star"`
	IsBanned                 bool       `json:"isBanned" db:"isBanned"`
	BanReason                *string    `json:"banReason,omitempty" db:"banReason"`
	BannedBy                 *string    `json:"bannedBy,omitempty" db:"bannedBy"`
	PostsCount               int        `json:"postsCount" db:"postsCount"`
	IsEmailVerified          bool       `json:"isEmailVerified" db:"isEmailVerified"`
	EmailVerificationToken   *string    `json:"-" db:"emailVerificationToken"`
	EmailVerificationExpires *time.Time `json:"-" db:"emailVerificationExpires"`
	PasswordResetToken       *string    `json:"-" db:"passwordResetToken"`
	PasswordResetExpires     *time.Time `json:"-" db:"passwordResetExpires"`
	LastLoginAt              *time.Time `json:"lastLoginAt,omitempty" db:"lastLoginAt"`
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
