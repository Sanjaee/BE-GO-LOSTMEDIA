package repositories

import (
	"errors"
	"time"

	"lostmediago/internal/models"
	"lostmediago/pkg/database"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(userId string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByGoogleID(googleId string) (*models.User, error)
	FindByEmailVerificationToken(token string) (*models.User, error)
	FindByPasswordResetToken(token string) (*models.User, error)
	Update(user *models.User) error
	UpdateLastLogin(userId string) error
	UpdateEmailVerification(userId string, isVerified bool) error
	UpdateEmailVerificationToken(userId, token string, expiresAt time.Time) error
	UpdatePasswordResetToken(userId, token string, expiresAt time.Time) error
	ClearPasswordResetToken(userId string) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByID(userId string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ? AND is_banned = ?", userId, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGoogleID(googleId string) (*models.User, error) {
	var user models.User
	err := r.db.Where("google_id = ?", googleId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmailVerificationToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email_verification_token = ? AND email_verification_expires > ?", token, time.Now()).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired verification token")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPasswordResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(userId string) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", userId).Update("last_login_at", time.Now()).Error
}

func (r *userRepository) UpdateEmailVerification(userId string, isVerified bool) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"is_email_verified":          isVerified,
		"email_verification_token":   nil,
		"email_verification_expires": nil,
	}).Error
}

func (r *userRepository) UpdateEmailVerificationToken(userId, token string, expiresAt time.Time) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"email_verification_token":   token,
		"email_verification_expires": expiresAt,
	}).Error
}

func (r *userRepository) UpdatePasswordResetToken(userId, token string, expiresAt time.Time) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"password_reset_token":   token,
		"password_reset_expires": expiresAt,
	}).Error
}

func (r *userRepository) ClearPasswordResetToken(userId string) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"password_reset_token":   nil,
		"password_reset_expires": nil,
	}).Error
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
