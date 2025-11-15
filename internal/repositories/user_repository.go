package repositories

import (
	"database/sql"
	"errors"
	"time"

	"lostmediago/internal/models"
	"lostmediago/pkg/database"
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
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			userId, googleId, username, email, password, profilePic, bio,
			role, star, isEmailVerified, emailVerificationToken, emailVerificationExpires,
			createdAt, updatedAt
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(
		query,
		user.UserId,
		database.NullString(user.GoogleId),
		user.Username,
		user.Email,
		database.NullString(user.Password),
		database.NullString(user.ProfilePic),
		database.NullString(user.Bio),
		user.Role,
		user.Star,
		user.IsEmailVerified,
		database.NullString(user.EmailVerificationToken),
		database.NullTime(user.EmailVerificationExpires),
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *userRepository) FindByID(userId string) (*models.User, error) {
	query := `SELECT * FROM users WHERE userId = $1 AND isBanned = false`
	user, err := r.scanUser(query, userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	user, err := r.scanUser(query, email)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	user, err := r.scanUser(query, username)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByGoogleID(googleId string) (*models.User, error) {
	query := `SELECT * FROM users WHERE googleId = $1`
	user, err := r.scanUser(query, googleId)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByEmailVerificationToken(token string) (*models.User, error) {
	query := `
		SELECT * FROM users 
		WHERE emailVerificationToken = $1 
		AND emailVerificationExpires > NOW()
	`
	user, err := r.scanUser(query, token)
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid or expired verification token")
	}
	return user, err
}

func (r *userRepository) FindByPasswordResetToken(token string) (*models.User, error) {
	query := `
		SELECT * FROM users 
		WHERE passwordResetToken = $1 
		AND passwordResetExpires > NOW()
	`
	user, err := r.scanUser(query, token)
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid or expired reset token")
	}
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET
			googleId = $2,
			username = $3,
			email = $4,
			password = $5,
			profilePic = $6,
			bio = $7,
			updatedAt = $8,
			role = $9,
			star = $10,
			isBanned = $11,
			banReason = $12,
			bannedBy = $13,
			postsCount = $14,
			isEmailVerified = $15,
			emailVerificationToken = $16,
			emailVerificationExpires = $17,
			passwordResetToken = $18,
			passwordResetExpires = $19,
			lastLoginAt = $20
		WHERE userId = $1
	`

	_, err := r.db.Exec(
		query,
		user.UserId,
		database.NullString(user.GoogleId),
		user.Username,
		user.Email,
		database.NullString(user.Password),
		database.NullString(user.ProfilePic),
		database.NullString(user.Bio),
		user.UpdatedAt,
		user.Role,
		user.Star,
		user.IsBanned,
		database.NullString(user.BanReason),
		database.NullString(user.BannedBy),
		user.PostsCount,
		user.IsEmailVerified,
		database.NullString(user.EmailVerificationToken),
		database.NullTime(user.EmailVerificationExpires),
		database.NullString(user.PasswordResetToken),
		database.NullTime(user.PasswordResetExpires),
		database.NullTime(user.LastLoginAt),
	)

	return err
}

func (r *userRepository) UpdateLastLogin(userId string) error {
	query := `UPDATE users SET lastLoginAt = NOW() WHERE userId = $1`
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *userRepository) UpdateEmailVerification(userId string, isVerified bool) error {
	query := `
		UPDATE users SET 
			isEmailVerified = $2,
			emailVerificationToken = NULL,
			emailVerificationExpires = NULL
		WHERE userId = $1
	`
	_, err := r.db.Exec(query, userId, isVerified)
	return err
}

func (r *userRepository) UpdateEmailVerificationToken(userId, token string, expiresAt time.Time) error {
	query := `
		UPDATE users SET 
			emailVerificationToken = $2,
			emailVerificationExpires = $3
		WHERE userId = $1
	`
	_, err := r.db.Exec(query, userId, token, expiresAt)
	return err
}

func (r *userRepository) UpdatePasswordResetToken(userId, token string, expiresAt time.Time) error {
	query := `
		UPDATE users SET 
			passwordResetToken = $2,
			passwordResetExpires = $3
		WHERE userId = $1
	`
	_, err := r.db.Exec(query, userId, token, expiresAt)
	return err
}

func (r *userRepository) ClearPasswordResetToken(userId string) error {
	query := `UPDATE users SET passwordResetToken = NULL, passwordResetExpires = NULL WHERE userId = $1`
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (r *userRepository) scanUser(query string, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	var (
		googleId, password, profilePic, bio, banReason, bannedBy sql.NullString
		emailVerificationToken, passwordResetToken               sql.NullString
		emailVerificationExpires, passwordResetExpires           sql.NullTime
		lastLoginAt                                              sql.NullTime
	)

	err := r.db.QueryRow(query, args...).Scan(
		&user.UserId,
		&googleId,
		&user.Username,
		&user.Email,
		&password,
		&profilePic,
		&bio,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.Role,
		&user.Star,
		&user.IsBanned,
		&banReason,
		&bannedBy,
		&user.PostsCount,
		&user.IsEmailVerified,
		&emailVerificationToken,
		&emailVerificationExpires,
		&passwordResetToken,
		&passwordResetExpires,
		&lastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	if googleId.Valid {
		user.GoogleId = &googleId.String
	}
	if password.Valid {
		user.Password = &password.String
	}
	if profilePic.Valid {
		user.ProfilePic = &profilePic.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if banReason.Valid {
		user.BanReason = &banReason.String
	}
	if bannedBy.Valid {
		user.BannedBy = &bannedBy.String
	}
	if emailVerificationToken.Valid {
		user.EmailVerificationToken = &emailVerificationToken.String
	}
	if emailVerificationExpires.Valid {
		user.EmailVerificationExpires = &emailVerificationExpires.Time
	}
	if passwordResetToken.Valid {
		user.PasswordResetToken = &passwordResetToken.String
	}
	if passwordResetExpires.Valid {
		user.PasswordResetExpires = &passwordResetExpires.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}
