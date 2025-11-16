package services

import (
	"context"
	"fmt"
	"log"
	"lostmediago/internal/models"
	"lostmediago/internal/repositories"
	"lostmediago/pkg/search"
)

type SearchService struct {
	postRepo repositories.PostRepository
}

func NewSearchService(postRepo repositories.PostRepository) *SearchService {
	return &SearchService{
		postRepo: postRepo,
	}
}

// IndexPost indexes a post in Bleve
func (s *SearchService) IndexPost(ctx context.Context, post *models.Post) error {
	// Only index published and not deleted posts
	if !post.IsPublished || post.IsDeleted {
		return nil
	}

	// Prepare document for indexing
	doc := map[string]interface{}{
		"postId":      post.PostId,
		"userId":      post.UserId,
		"title":       post.Title,
		"description": "",
		"content":     "",
		"category":    post.Category,
		"createdAt":   post.CreatedAt,
	}

	if post.Description != nil {
		doc["description"] = *post.Description
	}
	if post.Content != nil {
		doc["content"] = *post.Content
	}

	// Index the document
	if err := search.IndexDocument(post.PostId, doc); err != nil {
		log.Printf("[SEARCH ERROR] Failed to index post %s: %v", post.PostId, err)
		return err
	}

	log.Printf("[SEARCH] Indexed post: %s - %s", post.PostId, post.Title)
	return nil
}

// SearchPosts searches for posts using Bleve
func (s *SearchService) SearchPosts(ctx context.Context, query string, limit, offset int) ([]models.Post, int64, error) {
	log.Printf("[SEARCH] Searching for query: '%s' (limit: %d, offset: %d)", query, limit, offset)

	// Perform search
	postIds, total, err := search.Search(query, limit, offset)
	if err != nil {
		log.Printf("[SEARCH ERROR] Bleve search failed: %v", err)
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	log.Printf("[SEARCH] Bleve returned %d post IDs (total: %d)", len(postIds), total)

	if len(postIds) == 0 {
		log.Printf("[SEARCH] No post IDs found, returning empty result")
		return []models.Post{}, total, nil
	}

	// Fetch full post data from database
	posts, err := s.postRepo.FindByIDs(postIds)
	if err != nil {
		log.Printf("[SEARCH ERROR] Failed to fetch posts from DB: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}

	log.Printf("[SEARCH] Fetched %d posts from database (requested %d IDs)", len(posts), len(postIds))

	// Maintain search result order
	postMap := make(map[string]*models.Post)
	for i := range posts {
		postMap[posts[i].PostId] = &posts[i]
	}

	orderedPosts := make([]models.Post, 0, len(postIds))
	for _, postId := range postIds {
		if post, ok := postMap[postId]; ok {
			orderedPosts = append(orderedPosts, *post)
		} else {
			log.Printf("[SEARCH WARNING] Post ID %s from search not found in database", postId)
		}
	}

	log.Printf("[SEARCH] Returning %d ordered posts", len(orderedPosts))
	return orderedPosts, total, nil
}

// DeletePost removes a post from the search index
func (s *SearchService) DeletePost(ctx context.Context, postId string) error {
	return search.DeleteDocument(postId)
}

// ReindexAllPosts reindexes all published posts
func (s *SearchService) ReindexAllPosts(ctx context.Context) error {
	log.Println("[SEARCH] Starting reindex of all published posts...")

	// Fetch all published posts
	posts, err := s.postRepo.FindAllPublished(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	log.Printf("[SEARCH] Found %d published posts to index", len(posts))

	// Index each post
	successCount := 0
	errorCount := 0
	for i, post := range posts {
		if err := s.IndexPost(ctx, &post); err != nil {
			log.Printf("[SEARCH WARNING] Failed to index post %s (%d/%d): %v", post.PostId, i+1, len(posts), err)
			errorCount++
		} else {
			successCount++
			if (i+1)%10 == 0 {
				log.Printf("[SEARCH] Progress: %d/%d posts indexed", i+1, len(posts))
			}
		}
	}

	log.Printf("[SEARCH] Reindex completed: %d/%d posts indexed successfully (%d errors)", successCount, len(posts), errorCount)

	// Verify index count
	if idx := search.GetIndex(); idx != nil {
		docCount, err := idx.DocCount()
		if err == nil {
			log.Printf("[SEARCH] Index now contains %d documents", docCount)
		}
	}

	return nil
}
