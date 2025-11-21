package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// BroadcastRepository handles broadcast data operations
type BroadcastRepository struct {
	db *gorm.DB
}

// NewBroadcastRepository creates a new broadcast repository
func NewBroadcastRepository(db *gorm.DB) *BroadcastRepository {
	return &BroadcastRepository{db: db}
}

// CreateBroadcast creates a new broadcast
func (r *BroadcastRepository) CreateBroadcast(ctx context.Context, broadcast *models.Broadcast) error {
	return r.db.WithContext(ctx).Create(broadcast).Error
}

// GetBroadcastByID gets a broadcast by ID
func (r *BroadcastRepository) GetBroadcastByID(ctx context.Context, id string) (*models.Broadcast, error) {
	var broadcast models.Broadcast
	if err := r.db.WithContext(ctx).Preload("Audiences").Where("id = ?", id).First(&broadcast).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("broadcast not found")
		}
		return nil, err
	}
	return &broadcast, nil
}

// GetAllBroadcasts gets all broadcasts with filters
func (r *BroadcastRepository) GetAllBroadcasts(ctx context.Context, broadcastType *string, isPublished *bool, limit, offset int) ([]models.Broadcast, int64, error) {
	var broadcasts []models.Broadcast
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Broadcast{})

	if broadcastType != nil {
		query = query.Where("type = ?", *broadcastType)
	}
	if isPublished != nil {
		query = query.Where("is_published = ?", *isPublished)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Audiences").Order("created_at DESC").Limit(limit).Offset(offset).Find(&broadcasts).Error; err != nil {
		return nil, 0, err
	}

	return broadcasts, total, nil
}

// UpdateBroadcast updates a broadcast
func (r *BroadcastRepository) UpdateBroadcast(ctx context.Context, broadcast *models.Broadcast) error {
	return r.db.WithContext(ctx).Save(broadcast).Error
}

// DeleteBroadcast soft deletes a broadcast
func (r *BroadcastRepository) DeleteBroadcast(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Broadcast{}, "id = ?", id).Error
}

// SearchBroadcasts searches broadcasts by title or content
func (r *BroadcastRepository) SearchBroadcasts(ctx context.Context, query string, limit, offset int) ([]models.Broadcast, int64, error) {
	var broadcasts []models.Broadcast
	var total int64

	searchPattern := "%" + query + "%"
	dbQuery := r.db.WithContext(ctx).Model(&models.Broadcast{}).
		Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern).
		Where("is_published = ?", true)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Preload("Audiences").Order("created_at DESC").Limit(limit).Offset(offset).Find(&broadcasts).Error; err != nil {
		return nil, 0, err
	}

	return broadcasts, total, nil
}

// GetGeneralBroadcasts gets general broadcasts
func (r *BroadcastRepository) GetGeneralBroadcasts(ctx context.Context, limit, offset int) ([]models.Broadcast, int64, error) {
	return r.GetAllBroadcasts(ctx, func() *string { s := "general"; return &s }(), func() *bool { b := true; return &b }(), limit, offset)
}

// GetClassBroadcasts gets class-specific broadcasts
func (r *BroadcastRepository) GetClassBroadcasts(ctx context.Context, classID *string, limit, offset int) ([]models.Broadcast, int64, error) {
	var broadcasts []models.Broadcast
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Broadcast{}).
		Where("type = ? AND is_published = ?", "class", true)

	if classID != nil {
		query = query.Where("class_id = ?", *classID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Audiences").Order("created_at DESC").Limit(limit).Offset(offset).Find(&broadcasts).Error; err != nil {
		return nil, 0, err
	}

	return broadcasts, total, nil
}

// CreateAudience creates a broadcast audience
func (r *BroadcastRepository) CreateAudience(ctx context.Context, audience *models.BroadcastAudience) error {
	return r.db.WithContext(ctx).Create(audience).Error
}

// GetScheduledBroadcasts gets broadcasts scheduled for publishing
func (r *BroadcastRepository) GetScheduledBroadcasts(ctx context.Context, before time.Time) ([]models.Broadcast, error) {
	var broadcasts []models.Broadcast
	if err := r.db.WithContext(ctx).
		Where("scheduled_at IS NOT NULL AND scheduled_at <= ? AND is_published = ?", before, false).
		Find(&broadcasts).Error; err != nil {
		return nil, err
	}
	return broadcasts, nil
}

