package repositories

import (
	"context"
	"errors"
	"lostmediago/internal/models"
	"lostmediago/pkg/database"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(postId string) (*models.Post, error)
	FindByIDWithRelations(postId string, userId *string) (*models.Post, error)
	FindAll(limit, offset int, userId *string) ([]models.Post, int64, error)
	FindByUserID(userId string, limit, offset int) ([]models.Post, int64, error)
	FindByIDs(postIds []string) ([]models.Post, error)
	FindAllPublished(ctx context.Context) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(postId string) error
	IncrementViews(postId string) error
	GetUserPostsCount(userId string) (int, error)
	FindScheduledPosts() ([]models.Post, error)
	PublishScheduledPost(postId string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository() PostRepository {
	return &postRepository{
		db: database.DB,
	}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(postId string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("post_id = ? AND is_deleted = ?", postId, false).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindByIDWithRelations(postId string, userId *string) (*models.Post, error) {
	var post models.Post
	query := r.db.Where("post_id = ? AND is_deleted = ?", postId, false).
		Preload("User").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		})

	err := query.First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Check if user liked this post
	if userId != nil {
		var likeCount int64
		r.db.Model(&models.Like{}).
			Where("post_id = ? AND user_id = ?", postId, *userId).
			Count(&likeCount)
		// Note: We'll handle isLiked in service layer
	}

	return &post, nil
}

func (r *postRepository) FindAll(limit, offset int, userId *string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).
		Where("is_deleted = ? AND is_published = ?", false, true).
		Preload("User").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get posts
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindByUserID(userId string, limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).
		Where("user_id = ? AND is_deleted = ?", userId, false).
		Preload("User").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get posts
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(postId string) error {
	return r.db.Model(&models.Post{}).
		Where("post_id = ?", postId).
		Update("is_deleted", true).Error
}

func (r *postRepository) IncrementViews(postId string) error {
	return r.db.Model(&models.Post{}).
		Where("post_id = ?", postId).
		UpdateColumn("views_count", gorm.Expr("views_count + ?", 1)).Error
}

func (r *postRepository) GetUserPostsCount(userId string) (int, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Where("user_id = ? AND is_deleted = ?", userId, false).
		Count(&count).Error
	return int(count), err
}

func (r *postRepository) FindScheduledPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("is_scheduled = ? AND is_published = ? AND scheduled_at <= ? AND is_deleted = ?",
		true, false, "NOW()", false).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) PublishScheduledPost(postId string) error {
	return r.db.Model(&models.Post{}).
		Where("post_id = ?", postId).
		Updates(map[string]interface{}{
			"is_published": true,
			"is_scheduled": false,
		}).Error
}

func (r *postRepository) FindByIDs(postIds []string) ([]models.Post, error) {
	if len(postIds) == 0 {
		return []models.Post{}, nil
	}

	var posts []models.Post
	err := r.db.Where("post_id IN ? AND is_deleted = ? AND is_published = ?", postIds, false, true).
		Preload("User").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) FindAllPublished(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.WithContext(ctx).
		Where("is_deleted = ? AND is_published = ?", false, true).
		Preload("User").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Find(&posts).Error
	return posts, err
}
