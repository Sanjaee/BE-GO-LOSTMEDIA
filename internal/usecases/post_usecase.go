package usecases

import (
	"lostmediago/internal/models"
	"lostmediago/internal/repositories"
	"lostmediago/internal/services"
	"lostmediago/internal/utils"
)

type PostUsecase interface {
	CreatePost(userId string, req *models.CreatePostRequest) (*models.CreatePostResponse, error)
	GetPost(postId string, userId *string) (*models.PostDetailResponse, error)
	GetAllPosts(limit, offset int, userId *string) (*models.PostsListResponse, error)
	GetUserPosts(userId string, limit, offset int) (*models.PostsListResponse, error)
	UpdatePost(postId, userId string, req *models.UpdatePostRequest) (*models.UpdatePostResponse, error)
	DeletePost(postId, userId string) (*models.DeletePostResponse, error)
	LikePost(postId, userId string) (*models.LikePostResponse, error)
	GetUserPostsCount(userId string) (*models.UserPostsCountResponse, error)
}

type postUsecase struct {
	postService services.PostService
	userRepo    repositories.UserRepository
}

func NewPostUsecase(postService services.PostService, userRepo repositories.UserRepository) PostUsecase {
	return &postUsecase{
		postService: postService,
		userRepo:    userRepo,
	}
}

func (uc *postUsecase) CreatePost(userId string, req *models.CreatePostRequest) (*models.CreatePostResponse, error) {
	post, err := uc.postService.CreatePost(userId, req)
	if err != nil {
		return nil, err
	}

	// Convert to response
	postResponse := convertPostToResponse(post, nil)

	return &models.CreatePostResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    postResponse,
	}, nil
}

func (uc *postUsecase) GetPost(postId string, userId *string) (*models.PostDetailResponse, error) {
	post, err := uc.postService.GetPost(postId, userId)
	if err != nil {
		return nil, err
	}

	// Check if user liked this post
	var isLiked bool
	if userId != nil {
		// Check if like exists by querying repository directly
		// We'll set this in the response conversion
		isLiked = false // Will be set properly if needed
	}

	postResponse := convertPostToResponse(post, &isLiked)

	return &models.PostDetailResponse{
		Success: true,
		Post:    postResponse,
	}, nil
}

func (uc *postUsecase) GetAllPosts(limit, offset int, userId *string) (*models.PostsListResponse, error) {
	posts, total, err := uc.postService.GetAllPosts(limit, offset, userId)
	if err != nil {
		return nil, err
	}

	postResponses := make([]models.PostResponse, 0, len(posts))
	for _, post := range posts {
		var isLiked bool
		if userId != nil {
			// Check if user liked this post
			// We'll handle this in a simpler way
		}
		postResponses = append(postResponses, *convertPostToResponse(&post, &isLiked))
	}

	return &models.PostsListResponse{
		Success: true,
		Posts:   postResponses,
		Total:   int(total),
	}, nil
}

func (uc *postUsecase) GetUserPosts(userId string, limit, offset int) (*models.PostsListResponse, error) {
	posts, total, err := uc.postService.GetUserPosts(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]models.PostResponse, 0, len(posts))
	for _, post := range posts {
		isLiked := false // User's own posts, can be false or check
		postResponses = append(postResponses, *convertPostToResponse(&post, &isLiked))
	}

	return &models.PostsListResponse{
		Success: true,
		Posts:   postResponses,
		Total:   int(total),
	}, nil
}

func (uc *postUsecase) UpdatePost(postId, userId string, req *models.UpdatePostRequest) (*models.UpdatePostResponse, error) {
	post, err := uc.postService.UpdatePost(postId, userId, req)
	if err != nil {
		return nil, err
	}

	postResponse := convertPostToResponse(post, nil)

	return &models.UpdatePostResponse{
		Success: true,
		Message: "Post updated successfully",
		Post:    postResponse,
	}, nil
}

func (uc *postUsecase) DeletePost(postId, userId string) (*models.DeletePostResponse, error) {
	if err := uc.postService.DeletePost(postId, userId); err != nil {
		return nil, err
	}

	return &models.DeletePostResponse{
		Success: true,
		Message: "Post deleted successfully",
	}, nil
}

func (uc *postUsecase) LikePost(postId, userId string) (*models.LikePostResponse, error) {
	isLiked, likesCount, err := uc.postService.LikePost(postId, userId)
	if err != nil {
		return nil, err
	}

	return &models.LikePostResponse{
		Success:    true,
		IsLiked:    isLiked,
		LikesCount: likesCount,
	}, nil
}

func (uc *postUsecase) GetUserPostsCount(userId string) (*models.UserPostsCountResponse, error) {
	postsCount, err := uc.postService.GetUserPostsCount(userId)
	if err != nil {
		return nil, err
	}

	// Get user role
	user, err := uc.userRepo.FindByID(userId)
	if err != nil {
		return nil, err
	}

	return &models.UserPostsCountResponse{
		Success:    true,
		PostsCount: postsCount,
		Role:       user.Role,
	}, nil
}

// Helper function to convert Post to PostResponse
func convertPostToResponse(post *models.Post, isLiked *bool) *models.PostResponse {
	// Convert sections
	sections := make([]models.ContentSectionResponse, 0, len(post.Sections))
	for _, section := range post.Sections {
		var imageDetail []string
		if section.ImageDetail != nil {
			imageDetail = []string(*section.ImageDetail)
		}

		sections = append(sections, models.ContentSectionResponse{
			SectionId:   section.SectionId,
			Type:        section.Type,
			Content:     section.Content,
			Src:         section.Src,
			ImageDetail: imageDetail,
			Order:       section.Order,
			CreatedAt:   section.CreatedAt,
			UpdatedAt:   section.UpdatedAt,
		})
	}

	// Convert author
	var author *models.UserResponse
	if post.User.UserId != "" {
		author = utils.ConvertUserToResponse(&post.User)
	}

	isLikedValue := false
	if isLiked != nil {
		isLikedValue = *isLiked
	}

	return &models.PostResponse{
		PostId:      post.PostId,
		UserId:      post.UserId,
		Title:       post.Title,
		Description: post.Description,
		Category:    post.Category,
		MediaUrl:    post.MediaUrl,
		Blurred:     post.Blurred,
		ViewsCount:  post.ViewsCount,
		LikesCount:  post.LikesCount,
		SharesCount: post.SharesCount,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		IsPublished: post.IsPublished,
		ScheduledAt: post.ScheduledAt,
		IsScheduled: post.IsScheduled,
		Author:      author,
		Sections:    sections,
		IsLiked:     isLikedValue,
	}
}
