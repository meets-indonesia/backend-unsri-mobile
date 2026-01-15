package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"unsri-backend/internal/shared/models"

	"gorm.io/gorm"
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

// haversineDistance calculates the distance between two points using Haversine formula
// Returns distance in meters
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6371000 // Earth radius in meters

	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusMeters * c

	return distance
}

// CheckLocationInGeofence checks if location is within geofence using Haversine formula
func (r *LocationRepository) CheckLocationInGeofence(ctx context.Context, latitude, longitude float64) (*models.Geofence, error) {
	var geofences []models.Geofence
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&geofences).Error; err != nil {
		return nil, err
	}

	fmt.Printf("[CheckLocationInGeofence] Checking %d active geofences for location (%f, %f)\n", len(geofences), latitude, longitude)

	// Check if location is within any active geofence
	for _, geofence := range geofences {
		// Calculate distance using Haversine formula (in meters)
		distance := haversineDistance(
			geofence.Latitude,
			geofence.Longitude,
			latitude,
			longitude,
		)

		fmt.Printf("[CheckLocationInGeofence] ID=%s, Lat=%f, Long=%f, Radius=%f, Distance=%f\n",
			geofence.ID, geofence.Latitude, geofence.Longitude, geofence.Radius, distance)

		// Check if distance is within radius (radius is in meters)
		if distance <= geofence.Radius {
			return &geofence, nil
		}
	}

	return nil, errors.New("location not in any geofence")
}
