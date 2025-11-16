package repositories

import (
	"errors"
	"lostmediago/internal/models"
	"lostmediago/pkg/database"

	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(like *models.Like) error
	FindByID(likeId string) (*models.Like, error)
	FindByUserAndPost(userId, postId string) (*models.Like, error)
	FindByUserAndComment(userId, commentId string) (*models.Like, error)
	Delete(likeId string) error
	CountByPost(postId string) (int64, error)
	CountByComment(commentId string) (int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository() LikeRepository {
	return &likeRepository{
		db: database.DB,
	}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) FindByID(likeId string) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("like_id = ?", likeId).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like not found")
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) FindByUserAndPost(userId, postId string) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND post_id = ?", userId, postId).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like not found")
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) FindByUserAndComment(userId, commentId string) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND comment_id = ?", userId, commentId).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like not found")
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) Delete(likeId string) error {
	return r.db.Where("like_id = ?", likeId).Delete(&models.Like{}).Error
}

func (r *likeRepository) CountByPost(postId string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("post_id = ?", postId).Count(&count).Error
	return count, err
}

func (r *likeRepository) CountByComment(commentId string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("comment_id = ?", commentId).Count(&count).Error
	return count, err
}
