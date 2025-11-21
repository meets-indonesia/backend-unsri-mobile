package service

import (
	"context"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/broadcast/repository"
)

// BroadcastService handles broadcast business logic
type BroadcastService struct {
	repo *repository.BroadcastRepository
}

// NewBroadcastService creates a new broadcast service
func NewBroadcastService(repo *repository.BroadcastRepository) *BroadcastService {
	return &BroadcastService{repo: repo}
}

// CreateBroadcastRequest represents create broadcast request
type CreateBroadcastRequest struct {
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=general class campus"`
	Priority   string    `json:"priority,omitempty"`
	ClassID    *string   `json:"class_id,omitempty"`
	ScheduledAt *string  `json:"scheduled_at,omitempty"`
	ExpiresAt  *string   `json:"expires_at,omitempty"`
	Audiences  []AudienceRequest `json:"audiences,omitempty"`
}

// AudienceRequest represents audience target
type AudienceRequest struct {
	UserID *string `json:"user_id,omitempty"`
	Role   *string `json:"role,omitempty"`
	Prodi  *string `json:"prodi,omitempty"`
}

// CreateBroadcast creates a new broadcast
func (s *BroadcastService) CreateBroadcast(ctx context.Context, createdBy string, req CreateBroadcastRequest) (*models.Broadcast, error) {
	broadcast := &models.Broadcast{
		Title:       req.Title,
		Content:     req.Content,
		Type:        models.BroadcastType(req.Type),
		Priority:    "normal",
		CreatedBy:   createdBy,
		ClassID:     req.ClassID,
		IsPublished: false,
	}

	if req.Priority != "" {
		broadcast.Priority = req.Priority
	}

	if req.ScheduledAt != nil {
		scheduledAt, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid scheduled_at format, use RFC3339")
		}
		broadcast.ScheduledAt = &scheduledAt
	}

	if req.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid expires_at format, use RFC3339")
		}
		broadcast.ExpiresAt = &expiresAt
	}

	if err := s.repo.CreateBroadcast(ctx, broadcast); err != nil {
		return nil, apperrors.NewInternalError("failed to create broadcast", err)
	}

	// Create audiences
	for _, aud := range req.Audiences {
		audience := &models.BroadcastAudience{
			BroadcastID: broadcast.ID,
			UserID:      aud.UserID,
			Role:        aud.Role,
			Prodi:       aud.Prodi,
		}
		if err := s.repo.CreateAudience(ctx, audience); err != nil {
			// Log error but continue
			continue
		}
	}

	return s.repo.GetBroadcastByID(ctx, broadcast.ID)
}

// GetBroadcastByID gets a broadcast by ID
func (s *BroadcastService) GetBroadcastByID(ctx context.Context, id string) (*models.Broadcast, error) {
	return s.repo.GetBroadcastByID(ctx, id)
}

// GetBroadcastsRequest represents get broadcasts request
type GetBroadcastsRequest struct {
	Type        string `form:"type"`
	IsPublished *bool  `form:"is_published"`
	Page        int    `form:"page,default=1"`
	PerPage     int    `form:"per_page,default=20"`
}

// GetBroadcasts gets all broadcasts
func (s *BroadcastService) GetBroadcasts(ctx context.Context, req GetBroadcastsRequest) ([]models.Broadcast, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var typePtr *string
	if req.Type != "" {
		typePtr = &req.Type
	}

	return s.repo.GetAllBroadcasts(ctx, typePtr, req.IsPublished, perPage, (page-1)*perPage)
}

// UpdateBroadcastRequest represents update broadcast request
type UpdateBroadcastRequest struct {
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	Priority   *string `json:"priority,omitempty"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

// UpdateBroadcast updates a broadcast
func (s *BroadcastService) UpdateBroadcast(ctx context.Context, id string, req UpdateBroadcastRequest) (*models.Broadcast, error) {
	broadcast, err := s.repo.GetBroadcastByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("broadcast", id)
	}

	if req.Title != nil {
		broadcast.Title = *req.Title
	}
	if req.Content != nil {
		broadcast.Content = *req.Content
	}
	if req.Priority != nil {
		broadcast.Priority = *req.Priority
	}
	if req.IsPublished != nil {
		broadcast.IsPublished = *req.IsPublished
		if *req.IsPublished && broadcast.PublishedAt == nil {
			now := time.Now()
			broadcast.PublishedAt = &now
		}
	}

	if err := s.repo.UpdateBroadcast(ctx, broadcast); err != nil {
		return nil, apperrors.NewInternalError("failed to update broadcast", err)
	}

	return broadcast, nil
}

// DeleteBroadcast deletes a broadcast
func (s *BroadcastService) DeleteBroadcast(ctx context.Context, id string) error {
	_, err := s.repo.GetBroadcastByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("broadcast", id)
	}
	return s.repo.DeleteBroadcast(ctx, id)
}

// SearchBroadcasts searches broadcasts
func (s *BroadcastService) SearchBroadcasts(ctx context.Context, query string, page, perPage int) ([]models.Broadcast, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return s.repo.SearchBroadcasts(ctx, query, perPage, (page-1)*perPage)
}

// GetGeneralBroadcasts gets general broadcasts
func (s *BroadcastService) GetGeneralBroadcasts(ctx context.Context, page, perPage int) ([]models.Broadcast, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return s.repo.GetGeneralBroadcasts(ctx, perPage, (page-1)*perPage)
}

// GetClassBroadcasts gets class broadcasts
func (s *BroadcastService) GetClassBroadcasts(ctx context.Context, classID *string, page, perPage int) ([]models.Broadcast, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return s.repo.GetClassBroadcasts(ctx, classID, perPage, (page-1)*perPage)
}

// ScheduleBroadcastRequest represents schedule broadcast request
type ScheduleBroadcastRequest struct {
	ScheduledAt string `json:"scheduled_at" binding:"required"`
}

// ScheduleBroadcast schedules a broadcast
func (s *BroadcastService) ScheduleBroadcast(ctx context.Context, id string, req ScheduleBroadcastRequest) (*models.Broadcast, error) {
	broadcast, err := s.repo.GetBroadcastByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("broadcast", id)
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid scheduled_at format, use RFC3339")
	}

	broadcast.ScheduledAt = &scheduledAt

	if err := s.repo.UpdateBroadcast(ctx, broadcast); err != nil {
		return nil, apperrors.NewInternalError("failed to schedule broadcast", err)
	}

	return broadcast, nil
}

