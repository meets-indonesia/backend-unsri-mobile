package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// LocationRepository handles location data operations
type LocationRepository struct {
	db *gorm.DB
}

// NewLocationRepository creates a new location repository
func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// CreateGeofence creates a new geofence
func (r *LocationRepository) CreateGeofence(ctx context.Context, geofence *models.Geofence) error {
	return r.db.WithContext(ctx).Create(geofence).Error
}

// GetGeofenceByID gets a geofence by ID
func (r *LocationRepository) GetGeofenceByID(ctx context.Context, id string) (*models.Geofence, error) {
	var geofence models.Geofence
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&geofence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("geofence not found")
		}
		return nil, err
	}
	return &geofence, nil
}

// GetAllGeofences gets all active geofences
func (r *LocationRepository) GetAllGeofences(ctx context.Context) ([]models.Geofence, error) {
	var geofences []models.Geofence
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&geofences).Error; err != nil {
		return nil, err
	}
	return geofences, nil
}

// UpdateGeofence updates a geofence
func (r *LocationRepository) UpdateGeofence(ctx context.Context, geofence *models.Geofence) error {
	return r.db.WithContext(ctx).Save(geofence).Error
}

// DeleteGeofence soft deletes a geofence
func (r *LocationRepository) DeleteGeofence(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Geofence{}, "id = ?", id).Error
}

// CreateLocationHistory creates a location history record
func (r *LocationRepository) CreateLocationHistory(ctx context.Context, history *models.LocationHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetLocationHistory gets location history for a user
func (r *LocationRepository) GetLocationHistory(ctx context.Context, userID string, limit, offset int) ([]models.LocationHistory, int64, error) {
	var history []models.LocationHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.LocationHistory{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// GetCurrentTapInStatus gets current tap in status
func (r *LocationRepository) GetCurrentTapInStatus(ctx context.Context, userID string) (*models.LocationHistory, error) {
	var history models.LocationHistory
	today := time.Now().Format("2006-01-02")
	
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND DATE(created_at) = ?", userID, "tap_in", today).
		Order("created_at DESC").
		First(&history).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

// CheckLocationInGeofence checks if location is within geofence
func (r *LocationRepository) CheckLocationInGeofence(ctx context.Context, latitude, longitude float64) (*models.Geofence, error) {
	var geofences []models.Geofence
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&geofences).Error; err != nil {
		return nil, err
	}

	// Simple distance calculation (Haversine formula would be better)
	for _, geofence := range geofences {
		// Calculate distance (simplified)
		latDiff := latitude - geofence.Latitude
		lonDiff := longitude - geofence.Longitude
		distance := latDiff*latDiff + lonDiff*lonDiff
		
		// Convert radius to degrees (approximate: 1 degree ≈ 111km)
		radiusInDegrees := geofence.Radius / 111000.0
		
		if distance <= radiusInDegrees*radiusInDegrees {
			return &geofence, nil
		}
	}

	return nil, errors.New("location not in any geofence")
}

