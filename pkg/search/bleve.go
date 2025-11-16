package search

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

var (
	index bleve.Index
	mu    sync.RWMutex
)

// Connect initializes Bleve search index
func Connect() error {
	mu.Lock()
	defer mu.Unlock()

	// Get index path from environment or use default
	indexPath := os.Getenv("BLEVE_INDEX_PATH")
	if indexPath == "" {
		indexPath = "./data/bleve"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(indexPath, 0755); err != nil {
		return fmt.Errorf("failed to create index directory: %w", err)
	}

	indexPath = filepath.Join(indexPath, "posts.bleve")

	// Try to open existing index
	var err error
	index, err = bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// Create new index
		log.Println("[BLEVE] Creating new index...")
		mapping := buildIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		log.Println("[BLEVE] ✓ Index created successfully")
	} else if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	} else {
		log.Println("[BLEVE] ✓ Index opened successfully")
	}

	// Get document count
	docCount, err := index.DocCount()
	if err != nil {
		log.Printf("[BLEVE WARNING] Failed to get document count: %v", err)
	} else {
		log.Printf("[BLEVE] Index contains %d documents", docCount)
		if docCount == 0 {
			log.Printf("[BLEVE] WARNING: Index is empty! Posts need to be indexed.")
		}
	}

	return nil
}

// buildIndexMapping creates the index mapping for posts
func buildIndexMapping() mapping.IndexMapping {
	// Create a document mapping
	postMapping := bleve.NewDocumentMapping()

	// Title field - text, indexed, stored
	titleField := bleve.NewTextFieldMapping()
	titleField.Store = true
	titleField.Index = true
	titleField.Analyzer = "en"
	postMapping.AddFieldMappingsAt("title", titleField)

	// Description field - text, indexed, stored
	descField := bleve.NewTextFieldMapping()
	descField.Store = true
	descField.Index = true
	descField.Analyzer = "en"
	postMapping.AddFieldMappingsAt("description", descField)

	// Content field - text, indexed, stored
	contentField := bleve.NewTextFieldMapping()
	contentField.Store = true
	contentField.Index = true
	contentField.Analyzer = "en"
	postMapping.AddFieldMappingsAt("content", contentField)

	// Category field - text, indexed, stored
	categoryField := bleve.NewTextFieldMapping()
	categoryField.Store = true
	categoryField.Index = true
	categoryField.Analyzer = "keyword"
	postMapping.AddFieldMappingsAt("category", categoryField)

	// PostId field - keyword, stored (for retrieval)
	postIdField := bleve.NewKeywordFieldMapping()
	postIdField.Store = true
	postIdField.Index = true
	postMapping.AddFieldMappingsAt("postId", postIdField)

	// UserId field - keyword, indexed
	userIdField := bleve.NewKeywordFieldMapping()
	userIdField.Store = true
	userIdField.Index = true
	postMapping.AddFieldMappingsAt("userId", userIdField)

	// Metadata fields - not indexed, just stored
	createdAtField := bleve.NewDateTimeFieldMapping()
	createdAtField.Store = true
	createdAtField.Index = true
	postMapping.AddFieldMappingsAt("createdAt", createdAtField)

	// Create index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = postMapping
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}

// GetIndex returns the Bleve index
func GetIndex() bleve.Index {
	mu.RLock()
	defer mu.RUnlock()
	return index
}

// IndexDocument indexes a document
func IndexDocument(id string, data map[string]interface{}) error {
	mu.RLock()
	idx := index
	mu.RUnlock()

	if idx == nil {
		return fmt.Errorf("index not initialized")
	}

	if err := idx.Index(id, data); err != nil {
		return fmt.Errorf("failed to index document %s: %w", id, err)
	}

	return nil
}

// DeleteDocument removes a document from the index
func DeleteDocument(id string) error {
	mu.RLock()
	idx := index
	mu.RUnlock()

	if idx == nil {
		return fmt.Errorf("index not initialized")
	}

	return idx.Delete(id)
}

// Search performs a search query
func Search(query string, limit, offset int) ([]string, int64, error) {
	mu.RLock()
	idx := index
	mu.RUnlock()

	if idx == nil {
		return nil, 0, fmt.Errorf("index not initialized")
	}

	// If query is empty, return empty results
	if query == "" {
		return []string{}, 0, nil
	}

	// Create search query with multiple strategies for better partial matching
	// Use PrefixQuery for partial word matching (e.g., "fl" matches "flex")
	titlePrefixQuery := bleve.NewPrefixQuery(query)
	titlePrefixQuery.SetField("title")
	titlePrefixQuery.SetBoost(5.0)

	descPrefixQuery := bleve.NewPrefixQuery(query)
	descPrefixQuery.SetField("description")
	descPrefixQuery.SetBoost(3.0)

	contentPrefixQuery := bleve.NewPrefixQuery(query)
	contentPrefixQuery.SetField("content")
	contentPrefixQuery.SetBoost(2.0)

	categoryPrefixQuery := bleve.NewPrefixQuery(query)
	categoryPrefixQuery.SetField("category")
	categoryPrefixQuery.SetBoost(4.0)

	// Also create MatchQuery for exact word matching (higher boost for full words)
	titleMatchQuery := bleve.NewMatchQuery(query)
	titleMatchQuery.SetField("title")
	titleMatchQuery.SetBoost(6.0)
	titleMatchQuery.SetFuzziness(1) // Allow 1 character difference for typos

	descMatchQuery := bleve.NewMatchQuery(query)
	descMatchQuery.SetField("description")
	descMatchQuery.SetBoost(4.0)
	descMatchQuery.SetFuzziness(1)

	contentMatchQuery := bleve.NewMatchQuery(query)
	contentMatchQuery.SetField("content")
	contentMatchQuery.SetBoost(2.5)
	contentMatchQuery.SetFuzziness(1)

	categoryMatchQuery := bleve.NewMatchQuery(query)
	categoryMatchQuery.SetField("category")
	categoryMatchQuery.SetBoost(4.5)

	// Use QueryStringQuery for flexible matching (supports wildcards)
	// Format: "title:query* OR description:query* OR content:query*"
	queryString := fmt.Sprintf("title:%s* description:%s* content:%s* category:%s*", query, query, query, query)
	mainQuery := bleve.NewQueryStringQuery(queryString)
	mainQuery.SetBoost(3.0)

	// Combine all queries with OR (disjunction)
	disjunctQuery := bleve.NewDisjunctionQuery(
		titlePrefixQuery,
		descPrefixQuery,
		contentPrefixQuery,
		categoryPrefixQuery,
		titleMatchQuery,
		descMatchQuery,
		contentMatchQuery,
		categoryMatchQuery,
		mainQuery,
	)
	disjunctQuery.SetMin(1) // At least one query must match

	// Create search request
	searchRequest := bleve.NewSearchRequest(disjunctQuery)
	searchRequest.Size = limit
	searchRequest.From = offset
	// Sort by score descending first, then by date if scores are equal
	searchRequest.SortBy([]string{"-_score"})

	// Execute search
	searchResult, err := idx.Search(searchRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	log.Printf("[BLEVE] Search query '%s' returned %d results (total: %d)", query, len(searchResult.Hits), searchResult.Total)

	// Extract document IDs
	postIds := make([]string, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		postIds = append(postIds, hit.ID)
	}

	return postIds, int64(searchResult.Total), nil
}

// Close closes the index
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if index != nil {
		if err := index.Close(); err != nil {
			return fmt.Errorf("failed to close index: %w", err)
		}
		log.Println("[BLEVE] Index closed")
	}
	return nil
}
