package service

import (
	"context"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/search/repository"
)

// SearchService handles search business logic
type SearchService struct {
	repo *repository.SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(repo *repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// SearchRequest represents search request
type SearchRequest struct {
	Query   string                 `form:"q" binding:"required"`
	Type    string                 `form:"type"` // users, courses, schedules, global
	Role    *string                `form:"role"`
	Types   []string               `form:"types"` // for global search
	Filters map[string]interface{} `form:"filters"`
	Page    int                    `form:"page,default=1"`
	PerPage int                    `form:"per_page,default=20"`
}

// SearchResponse represents search response
type SearchResponse struct {
	Type    string      `json:"type"`
	Query   string      `json:"query"`
	Results interface{} `json:"results"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// GlobalSearchResponse represents global search response
type GlobalSearchResponse struct {
	Query   string                   `json:"query"`
	Results map[string]interface{}   `json:"results"`
}

// Search performs search based on type
func (s *SearchService) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.Query == "" {
		return nil, apperrors.NewValidationError("search query is required")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var results interface{}
	var total int64
	var err error

	switch req.Type {
	case "users":
		results, total, err = s.repo.SearchUsers(ctx, req.Query, req.Role, perPage, offset)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to search users", err)
		}

	case "courses":
		results, total, err = s.repo.SearchCourses(ctx, req.Query, perPage, offset)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to search courses", err)
		}

	case "schedules":
		var userID *string
		if userIDStr, ok := req.Filters["user_id"].(string); ok {
			userID = &userIDStr
		}
		results, total, err = s.repo.SearchSchedules(ctx, req.Query, userID, perPage, offset)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to search schedules", err)
		}

	default:
		// Type-specific search with filters
		results, total, err = s.repo.SearchByType(ctx, req.Type, req.Query, req.Filters, perPage, offset)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to search", err)
		}
	}

	return &SearchResponse{
		Type:    req.Type,
		Query:   req.Query,
		Results: results,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// GlobalSearch performs global search across multiple types
func (s *SearchService) GlobalSearch(ctx context.Context, query string, types []string, limit int) (*GlobalSearchResponse, error) {
	if query == "" {
		return nil, apperrors.NewValidationError("search query is required")
	}

	if len(types) == 0 {
		// Default types if not specified
		types = []string{"users", "courses", "schedules", "broadcasts"}
	}

	if limit < 1 {
		limit = 10
	}

	results, err := s.repo.SearchGlobal(ctx, query, types, limit)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to perform global search", err)
	}

	return &GlobalSearchResponse{
		Query:   query,
		Results: results,
	}, nil
}

