package handlers

import (
	"log"
	"lostmediago/internal/services"
	"lostmediago/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *services.SearchService
	postUsecase   usecases.PostUsecase
}

func NewSearchHandler(searchService *services.SearchService, postUsecase usecases.PostUsecase) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		postUsecase:   postUsecase,
	}
}

// SearchPosts handles POST search requests
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	var req struct {
		Query  string `json:"q" binding:"required"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}

	// Try to bind from JSON body first
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails, try query parameters
		req.Query = c.Query("q")
		if req.Query == "" {
			log.Printf("[SEARCH] Empty query received")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query parameter 'q' is required",
			})
			return
		}
		log.Printf("[SEARCH] Query from URL params: '%s'", req.Query)

		// Parse limit and offset from query params
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		var parseErr error
		req.Limit, parseErr = strconv.Atoi(limitStr)
		if parseErr != nil || req.Limit <= 0 {
			req.Limit = 10
		}
		if req.Limit > 100 {
			req.Limit = 100
		}

		req.Offset, parseErr = strconv.Atoi(offsetStr)
		if parseErr != nil || req.Offset < 0 {
			req.Offset = 0
		}
	} else {
		log.Printf("[SEARCH] Query from JSON body: '%s'", req.Query)
	}

	// Set defaults if not provided
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Perform search
	log.Printf("[SEARCH HANDLER] Processing search request - Query: '%s', Limit: %d, Offset: %d", req.Query, req.Limit, req.Offset)

	posts, total, err := h.searchService.SearchPosts(c.Request.Context(), req.Query, req.Limit, req.Offset)
	if err != nil {
		log.Printf("[SEARCH ERROR] Search failed for query '%s': %v", req.Query, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Search failed",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[SEARCH HANDLER] Query '%s' returned %d posts (total: %d)", req.Query, len(posts), total)

	// Convert to response format
	postResponses := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		// Get user info
		var author gin.H
		if post.User.UserId != "" {
			author = gin.H{
				"userId":     post.User.UserId,
				"username":   post.User.Username,
				"profilePic": post.User.ProfilePic,
			}
		}

		// Convert sections
		sections := make([]gin.H, 0, len(post.Sections))
		for _, section := range post.Sections {
			var imageDetail []string
			if section.ImageDetail != nil {
				imageDetail = []string(*section.ImageDetail)
			}
			sections = append(sections, gin.H{
				"sectionId":   section.SectionId,
				"type":        section.Type,
				"content":     section.Content,
				"src":         section.Src,
				"imageDetail": imageDetail,
				"order":       section.Order,
			})
		}

		postResponse := gin.H{
			"postId":      post.PostId,
			"userId":      post.UserId,
			"title":       post.Title,
			"description": post.Description,
			"category":    post.Category,
			"mediaUrl":    post.MediaUrl,
			"blurred":     post.Blurred,
			"viewsCount":  post.ViewsCount,
			"likesCount":  post.LikesCount,
			"sharesCount": post.SharesCount,
			"createdAt":   post.CreatedAt,
			"updatedAt":   post.UpdatedAt,
			"isPublished": post.IsPublished,
			"author":      author,
			"sections":    sections,
		}
		postResponses = append(postResponses, postResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"posts":  postResponses,
			"total":  total,
			"limit":  req.Limit,
			"offset": req.Offset,
		},
	})
}
