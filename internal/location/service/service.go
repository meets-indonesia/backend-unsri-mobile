package service

import (
	"context"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/location/repository"
)

// LocationService handles location business logic
type LocationService struct {
	repo *repository.LocationRepository
}

// NewLocationService creates a new location service
func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

// TapInRequest represents tap in request
type TapInRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// TapIn performs tap in with location validation
func (s *LocationService) TapIn(ctx context.Context, userID string, req TapInRequest) (*models.LocationHistory, error) {
	// Check if already tapped in today
	current, err := s.repo.GetCurrentTapInStatus(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check tap in status", err)
	}

	if current != nil {
		return nil, apperrors.NewConflictError("already tapped in today")
	}

	// Check if location is within geofence
	geofence, err := s.repo.CheckLocationInGeofence(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, apperrors.NewBadRequestError("location not within allowed area")
	}

	history := &models.LocationHistory{
		UserID:     userID,
		Type:       "tap_in",
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		GeofenceID: &geofence.ID,
		IsValid:    true,
	}

	if err := s.repo.CreateLocationHistory(ctx, history); err != nil {
		return nil, apperrors.NewInternalError("failed to record tap in", err)
	}

	return history, nil
}

// TapOutRequest represents tap out request
type TapOutRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// TapOut performs tap out
func (s *LocationService) TapOut(ctx context.Context, userID string, req TapOutRequest) (*models.LocationHistory, error) {
	// Check if tapped in today
	current, err := s.repo.GetCurrentTapInStatus(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check tap in status", err)
	}

	if current == nil {
		return nil, apperrors.NewBadRequestError("no active tap in found")
	}

	history := &models.LocationHistory{
		UserID:    userID,
		Type:      "tap_out",
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsValid:   true,
	}

	if err := s.repo.CreateLocationHistory(ctx, history); err != nil {
		return nil, apperrors.NewInternalError("failed to record tap out", err)
	}

	return history, nil
}

// GetCheckInStatus gets current check-in status
func (s *LocationService) GetCheckInStatus(ctx context.Context, userID string) (*models.LocationHistory, error) {
	return s.repo.GetCurrentTapInStatus(ctx, userID)
}

// GetLocationHistoryRequest represents get location history request
type GetLocationHistoryRequest struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

// GetLocationHistory gets location history
func (s *LocationService) GetLocationHistory(ctx context.Context, userID string, req GetLocationHistoryRequest) ([]models.LocationHistory, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.GetLocationHistory(ctx, userID, perPage, (page-1)*perPage)
}

// GetGeofences gets all geofences
func (s *LocationService) GetGeofences(ctx context.Context) ([]models.Geofence, error) {
	return s.repo.GetAllGeofences(ctx)
}

// ValidateLocationRequest represents validate location request
type ValidateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// ValidateLocation validates if location is within geofence
func (s *LocationService) ValidateLocation(ctx context.Context, req ValidateLocationRequest) (*models.Geofence, error) {
	geofence, err := s.repo.CheckLocationInGeofence(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, apperrors.NewBadRequestError("location not within allowed area")
	}
	return geofence, nil
}

// CreateGeofenceRequest represents create geofence request
type CreateGeofenceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description,omitempty"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Radius      float64 `json:"radius" binding:"required"` // in meters
}

// CreateGeofence creates a new geofence
func (s *LocationService) CreateGeofence(ctx context.Context, req CreateGeofenceRequest) (*models.Geofence, error) {
	geofence := &models.Geofence{
		Name:        req.Name,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		IsActive:    true,
	}

	if err := s.repo.CreateGeofence(ctx, geofence); err != nil {
		return nil, apperrors.NewInternalError("failed to create geofence", err)
	}

	return geofence, nil
}

