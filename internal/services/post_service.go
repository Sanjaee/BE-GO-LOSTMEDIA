package services

import (
	"context"
	"errors"
	"lostmediago/internal/models"
	"lostmediago/internal/repositories"
	"lostmediago/pkg/database"
	"time"

	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(userId string, req *models.CreatePostRequest) (*models.Post, error)
	GetPost(postId string, userId *string) (*models.Post, error)
	GetAllPosts(limit, offset int, userId *string) ([]models.Post, int64, error)
	GetUserPosts(userId string, limit, offset int) ([]models.Post, int64, error)
	UpdatePost(postId, userId string, req *models.UpdatePostRequest) (*models.Post, error)
	DeletePost(postId, userId string) error
	LikePost(postId, userId string) (bool, int, error) // returns isLiked, likesCount, error
	GetUserPostsCount(userId string) (int, error)
	IncrementViews(postId string) error
	PublishScheduledPosts() error
}

type postService struct {
	postRepo      repositories.PostRepository
	userRepo      repositories.UserRepository
	likeRepo      repositories.LikeRepository
	searchService *SearchService
}

func NewPostService(postRepo repositories.PostRepository, userRepo repositories.UserRepository, likeRepo repositories.LikeRepository) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
		likeRepo: likeRepo,
	}
}

func NewPostServiceWithSearch(postRepo repositories.PostRepository, userRepo repositories.UserRepository, likeRepo repositories.LikeRepository, searchService *SearchService) PostService {
	return &postService{
		postRepo:      postRepo,
		userRepo:      userRepo,
		likeRepo:      likeRepo,
		searchService: searchService,
	}
}

func (s *postService) CreatePost(userId string, req *models.CreatePostRequest) (*models.Post, error) {
	// Check if user exists
	if _, err := s.userRepo.FindByID(userId); err != nil {
		return nil, errors.New("user not found")
	}

	// Parse scheduled date if provided
	var scheduledAt *time.Time
	if req.IsScheduled && req.ScheduledAt != nil {
		parsed, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, errors.New("invalid scheduled date format")
		}
		scheduledAt = &parsed
	}

	// Determine if post should be published immediately
	// For now, publish immediately unless scheduled
	isPublished := !req.IsScheduled

	// Create post
	post := &models.Post{
		UserId:      userId,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		MediaUrl:    req.MediaUrl,
		Blurred:     req.Blurred,
		IsPublished: isPublished,
		IsScheduled: req.IsScheduled,
		ScheduledAt: scheduledAt,
	}

	// Create post first
	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// Create content sections
	if len(req.Sections) > 0 {
		sections := make([]models.ContentSection, 0, len(req.Sections))
		for _, sectionInput := range req.Sections {
			var imageDetail *models.ImageDetailArray
			if len(sectionInput.ImageDetail) > 0 {
				imgArray := models.ImageDetailArray(sectionInput.ImageDetail)
				imageDetail = &imgArray
			}

			section := models.ContentSection{
				Type:        sectionInput.Type,
				Content:     sectionInput.Content,
				Src:         sectionInput.Src,
				ImageDetail: imageDetail,
				Order:       sectionInput.Order,
				PostId:      post.PostId,
			}
			sections = append(sections, section)
		}

		// Save sections
		if err := database.DB.Create(&sections).Error; err != nil {
			return nil, err
		}
		post.Sections = sections
	}

	// Index post in BleveSearch if published
	if s.searchService != nil && post.IsPublished {
		go func() {
			if err := s.searchService.IndexPost(context.Background(), post); err != nil {
				// Log error but don't fail the request
			}
		}()
	}

	return post, nil
}

func (s *postService) GetPost(postId string, userId *string) (*models.Post, error) {
	post, err := s.postRepo.FindByIDWithRelations(postId, userId)
	if err != nil {
		return nil, err
	}

	// Increment views
	go s.postRepo.IncrementViews(postId)

	return post, nil
}

func (s *postService) GetAllPosts(limit, offset int, userId *string) ([]models.Post, int64, error) {
	return s.postRepo.FindAll(limit, offset, userId)
}

func (s *postService) GetUserPosts(userId string, limit, offset int) ([]models.Post, int64, error) {
	return s.postRepo.FindByUserID(userId, limit, offset)
}

func (s *postService) UpdatePost(postId, userId string, req *models.UpdatePostRequest) (*models.Post, error) {
	// Get existing post
	post, err := s.postRepo.FindByID(postId)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if post.UserId != userId {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	// Update fields
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Description != nil {
		post.Description = req.Description
	}
	if req.Category != nil {
		post.Category = *req.Category
	}
	if req.MediaUrl != nil {
		post.MediaUrl = req.MediaUrl
	}
	if req.Blurred != nil {
		post.Blurred = *req.Blurred
	}
	if req.IsScheduled != nil {
		post.IsScheduled = *req.IsScheduled
	}
	if req.ScheduledAt != nil {
		parsed, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			post.ScheduledAt = &parsed
		}
	}

	// Update post
	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	// Update sections if provided
	if req.Sections != nil {
		// Delete existing sections
		database.DB.Where("post_id = ?", postId).Delete(&models.ContentSection{})

		// Create new sections
		if len(req.Sections) > 0 {
			sections := make([]models.ContentSection, 0, len(req.Sections))
			for _, sectionInput := range req.Sections {
				var imageDetail *models.ImageDetailArray
				if len(sectionInput.ImageDetail) > 0 {
					imgArray := models.ImageDetailArray(sectionInput.ImageDetail)
					imageDetail = &imgArray
				}

				section := models.ContentSection{
					Type:        sectionInput.Type,
					Content:     sectionInput.Content,
					Src:         sectionInput.Src,
					ImageDetail: imageDetail,
					Order:       sectionInput.Order,
					PostId:      postId,
				}
				sections = append(sections, section)
			}

			if err := database.DB.Create(&sections).Error; err != nil {
				return nil, err
			}
		}
	}

	// Reload with relations
	updatedPost, err := s.postRepo.FindByIDWithRelations(postId, &userId)
	if err != nil {
		return nil, err
	}

	// Update index in BleveSearch if published
	if s.searchService != nil && updatedPost.IsPublished {
		go func() {
			if err := s.searchService.IndexPost(context.Background(), updatedPost); err != nil {
				// Log error but don't fail the request
			}
		}()
	}

	return updatedPost, nil
}

func (s *postService) DeletePost(postId, userId string) error {
	// Get post
	post, err := s.postRepo.FindByID(postId)
	if err != nil {
		return err
	}

	// Check ownership
	if post.UserId != userId {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	// Soft delete
	if err := s.postRepo.Delete(postId); err != nil {
		return err
	}

	// Remove from search index
	if s.searchService != nil {
		go func() {
			if err := s.searchService.DeletePost(context.Background(), postId); err != nil {
				// Log error but don't fail the request
			}
		}()
	}

	return nil
}

func (s *postService) LikePost(postId, userId string) (bool, int, error) {
	// Check if post exists
	if _, err := s.postRepo.FindByID(postId); err != nil {
		return false, 0, err
	}

	// Check if already liked
	existingLike, err := s.likeRepo.FindByUserAndPost(userId, postId)
	isLiked := err == nil && existingLike != nil

	if isLiked {
		// Unlike: delete like
		if err := s.likeRepo.Delete(existingLike.LikeId); err != nil {
			return false, 0, err
		}
		// Decrement likes count
		database.DB.Model(&models.Post{}).
			Where("post_id = ?", postId).
			UpdateColumn("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)"))

		// Get updated count
		updatedPost, _ := s.postRepo.FindByID(postId)
		return false, updatedPost.LikesCount, nil
	} else {
		// Like: create like
		like := &models.Like{
			UserId:   userId,
			PostId:   &postId,
			LikeType: "post",
		}
		if err := s.likeRepo.Create(like); err != nil {
			return false, 0, err
		}
		// Increment likes count
		database.DB.Model(&models.Post{}).
			Where("post_id = ?", postId).
			UpdateColumn("likes_count", gorm.Expr("likes_count + 1"))

		// Get updated count
		updatedPost, _ := s.postRepo.FindByID(postId)
		return true, updatedPost.LikesCount, nil
	}
}

func (s *postService) GetUserPostsCount(userId string) (int, error) {
	return s.postRepo.GetUserPostsCount(userId)
}

func (s *postService) IncrementViews(postId string) error {
	return s.postRepo.IncrementViews(postId)
}

func (s *postService) PublishScheduledPosts() error {
	posts, err := s.postRepo.FindScheduledPosts()
	if err != nil {
		return err
	}

	for _, post := range posts {
		if err := s.postRepo.PublishScheduledPost(post.PostId); err != nil {
			// Log error but continue with other posts
			continue
		}
	}

	return nil
}
