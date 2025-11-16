package models

import "time"

// CreatePostRequest represents request to create a post
type CreatePostRequest struct {
	Title       string                `json:"title" binding:"required,min=3"`
	Description *string               `json:"description,omitempty"`
	Category    string                `json:"category" binding:"required"`
	MediaUrl    *string               `json:"mediaUrl,omitempty"`
	Blurred     bool                  `json:"blurred"`
	Sections    []ContentSectionInput `json:"sections,omitempty"`
	ScheduledAt *string               `json:"scheduledAt,omitempty"` // ISO 8601 format
	IsScheduled bool                  `json:"isScheduled"`
}

// UpdatePostRequest represents request to update a post
type UpdatePostRequest struct {
	Title       *string               `json:"title,omitempty"`
	Description *string               `json:"description,omitempty"`
	Category    *string               `json:"category,omitempty"`
	MediaUrl    *string               `json:"mediaUrl,omitempty"`
	Blurred     *bool                 `json:"blurred,omitempty"`
	Sections    []ContentSectionInput `json:"sections,omitempty"`
	ScheduledAt *string               `json:"scheduledAt,omitempty"`
	IsScheduled *bool                 `json:"isScheduled,omitempty"`
}

// ContentSectionInput represents input for content section
type ContentSectionInput struct {
	Type        string   `json:"type" binding:"required,oneof=image code video link html"`
	Content     *string  `json:"content,omitempty"`
	Src         *string  `json:"src,omitempty"`
	ImageDetail []string `json:"imageDetail,omitempty"` // Array of image URLs
	Order       int      `json:"order"`
}

// PostResponse represents post response with relations
type PostResponse struct {
	PostId      string                   `json:"postId"`
	UserId      string                   `json:"userId"`
	Title       string                   `json:"title"`
	Description *string                  `json:"description,omitempty"`
	Category    string                   `json:"category"`
	MediaUrl    *string                  `json:"mediaUrl,omitempty"`
	Blurred     bool                     `json:"blurred"`
	ViewsCount  int                      `json:"viewsCount"`
	LikesCount  int                      `json:"likesCount"`
	SharesCount int                      `json:"sharesCount"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
	IsPublished bool                     `json:"isPublished"`
	ScheduledAt *time.Time               `json:"scheduledAt,omitempty"`
	IsScheduled bool                     `json:"isScheduled"`
	Author      *UserResponse            `json:"author,omitempty"`
	Sections    []ContentSectionResponse `json:"sections,omitempty"`
	IsLiked     bool                     `json:"isLiked,omitempty"`
}

// ContentSectionResponse represents content section in response
type ContentSectionResponse struct {
	SectionId   string    `json:"sectionId"`
	Type        string    `json:"type"`
	Content     *string   `json:"content,omitempty"`
	Src         *string   `json:"src,omitempty"`
	ImageDetail []string  `json:"imageDetail,omitempty"` // Array of image URLs
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// PostsListResponse represents list of posts response
type PostsListResponse struct {
	Success bool           `json:"success"`
	Posts   []PostResponse `json:"posts"`
	Total   int            `json:"total,omitempty"`
}

// PostDetailResponse represents single post response
type PostDetailResponse struct {
	Success bool          `json:"success"`
	Post    *PostResponse `json:"post,omitempty"`
}

// CreatePostResponse represents create post response
type CreatePostResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Post    *PostResponse `json:"post,omitempty"`
}

// UpdatePostResponse represents update post response
type UpdatePostResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Post    *PostResponse `json:"post,omitempty"`
}

// DeletePostResponse represents delete post response
type DeletePostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LikePostResponse represents like post response
type LikePostResponse struct {
	Success    bool `json:"success"`
	IsLiked    bool `json:"isLiked"`
	LikesCount int  `json:"likesCount"`
}

// UserPostsCountResponse represents user posts count response
type UserPostsCountResponse struct {
	Success    bool   `json:"success"`
	PostsCount int    `json:"postsCount"`
	Role       string `json:"role"`
}
