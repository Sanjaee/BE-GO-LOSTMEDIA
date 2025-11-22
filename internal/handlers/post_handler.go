package handlers

import (
	"net/http"
	"strconv"

	"lostmediago/internal/models"
	"lostmediago/internal/usecases"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUsecase usecases.PostUsecase
}

func NewPostHandler(postUsecase usecases.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

// CreatePost handles post creation
// @Summary Create a new post
// @Description Create a new post with content sections
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePostRequest true "Create post request"
// @Success 201 {object} Response{data=models.CreatePostResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "User ID not found in context",
			},
		})
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.postUsecase.CreatePost(userId.(string), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if errorMessage == "user not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Data: response,
	})
}

// GetPost handles getting a single post
// @Summary Get a post by ID
// @Description Get post details by post ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} Response{data=models.PostDetailResponse}
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "post ID is required",
			},
		})
		return
	}

	var userId *string
	if uid, exists := c.Get("userId"); exists {
		uidStr := uid.(string)
		userId = &uidStr
	}

	response, err := h.postUsecase.GetPost(postId, userId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if errorMessage == "post not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// GetAllPosts handles getting all posts
// @Summary Get all posts
// @Description Get list of all published posts
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} Response{data=models.PostsListResponse}
// @Router /api/v1/posts [get]
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var viewerUserId *string
	if uid, exists := c.Get("userId"); exists {
		uidStr := uid.(string)
		viewerUserId = &uidStr
	}

	// Optional explicit ownerId query (?userId=...) to filter by post owner
	ownerId := c.Query("userId")
	if ownerId != "" {
		// If user is authenticated, always use userId from JWT token (internal UUID)
		// This is because frontend may send Google ID but database uses UUID
		// Query param is ignored for authenticated users for security and correctness
		if viewerUserId != nil {
			// Use internal userId from JWT token (UUID) instead of query param
			response, err := h.postUsecase.GetUserPosts(*viewerUserId, limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: ErrorDetail{
						Code:    "INTERNAL_ERROR",
						Message: err.Error(),
					},
				})
				return
			}

			c.JSON(http.StatusOK, Response{
				Data: response,
			})
			return
		}

		// Unauthenticated user with query param - allow but may return empty
		// Note: This might not work if ownerId is Google ID and database uses UUID
		response, err := h.postUsecase.GetUserPosts(ownerId, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Data: response,
		})
		return
	}

	// Default: get all published posts (optionally personalized by viewer)
	response, err := h.postUsecase.GetAllPosts(limit, offset, viewerUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// UpdatePost handles post update
// @Summary Update a post
// @Description Update an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body models.UpdatePostRequest true "Update post request"
// @Success 200 {object} Response{data=models.UpdatePostResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "post ID is required",
			},
		})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "User ID not found in context",
			},
		})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.postUsecase.UpdatePost(postId, userId.(string), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if errorMessage == "post not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		} else if errorMessage == "unauthorized: you can only update your own posts" {
			statusCode = http.StatusForbidden
			errorCode = "FORBIDDEN"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// DeletePost handles post deletion
// @Summary Delete a post
// @Description Delete an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} Response{data=models.DeletePostResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "post ID is required",
			},
		})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "User ID not found in context",
			},
		})
		return
	}

	response, err := h.postUsecase.DeletePost(postId, userId.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if errorMessage == "post not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		} else if errorMessage == "unauthorized: you can only delete your own posts" {
			statusCode = http.StatusForbidden
			errorCode = "FORBIDDEN"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// LikePost handles post like/unlike
// @Summary Like or unlike a post
// @Description Toggle like status for a post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} Response{data=models.LikePostResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "post ID is required",
			},
		})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "User ID not found in context",
			},
		})
		return
	}

	response, err := h.postUsecase.LikePost(postId, userId.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_ERROR"
		errorMessage := err.Error()

		if errorMessage == "post not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: ErrorDetail{
				Code:    errorCode,
				Message: errorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}

// GetUserPostsCount handles getting user posts count
// @Summary Get user posts count
// @Description Get posts count and role for current user
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response{data=models.UserPostsCountResponse}
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/posts/user/posts-count [get]
func (h *PostHandler) GetUserPostsCount(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "User ID not found in context",
			},
		})
		return
	}

	response, err := h.postUsecase.GetUserPostsCount(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: response,
	})
}
