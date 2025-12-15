package service

import (
	"context"
	"time"

	"unsri-backend/internal/access/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// AccessService handles access control business logic
type AccessService struct {
	repo *repository.AccessRepository
}

// NewAccessService creates a new access service
func NewAccessService(repo *repository.AccessRepository) *AccessService {
	return &AccessService{repo: repo}
}

// ValidateQRRequest represents validate QR request
type ValidateQRRequest struct {
	QRToken string `json:"qr_token" binding:"required"`
	GateID  string `json:"gate_id" binding:"required"`
}

// ValidateQRResponse represents QR validation response
type ValidateQRResponse struct {
	Valid   bool   `json:"valid"`
	Allowed bool   `json:"allowed"`
	UserID  string `json:"user_id,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// ValidateAccessQR validates access QR code for gate
func (s *AccessService) ValidateAccessQR(ctx context.Context, req ValidateQRRequest) (*ValidateQRResponse, error) {
	// Get user from QR token (this should call QR service or user service)
	// For now, we'll assume the QR token contains user info
	// In production, this should validate with QR service first

	// Check permission
	permission, err := s.repo.GetAccessPermission(ctx, "", req.GateID) // userID should come from QR validation
	if err != nil {
		// Log access attempt
		if err := s.repo.CreateAccessLog(ctx, &models.AccessLog{
			GateID:     req.GateID,
			AccessType: "entry",
			IsAllowed:  false,
			Reason:     "permission not found",
		}); err != nil {
			// Log error but continue
			_ = err
		}

		return &ValidateQRResponse{
			Valid:   true,
			Allowed: false,
			Reason:  "no permission for this gate",
		}, nil
	}

	if !permission.IsAllowed {
		if err := s.repo.CreateAccessLog(ctx, &models.AccessLog{
			UserID:     permission.UserID,
			GateID:     req.GateID,
			AccessType: "entry",
			IsAllowed:  false,
			Reason:     "permission denied",
		}); err != nil {
			// Log error but continue
			_ = err
		}

		return &ValidateQRResponse{
			Valid:   true,
			Allowed: false,
			Reason:  "access denied",
		}, nil
	}

	// Log successful access
	if err := s.repo.CreateAccessLog(ctx, &models.AccessLog{
		UserID:     permission.UserID,
		GateID:     req.GateID,
		AccessType: "entry",
		IsAllowed:  true,
	}); err != nil {
		// Log error but continue
		_ = err
	}

	return &ValidateQRResponse{
		Valid:   true,
		Allowed: true,
		UserID:  permission.UserID,
	}, nil
}

// GetAccessHistoryRequest represents get access history request
type GetAccessHistoryRequest struct {
	UserID  *string `form:"user_id"`
	GateID  *string `form:"gate_id"`
	Page    int     `form:"page,default=1"`
	PerPage int     `form:"per_page,default=20"`
}

// GetAccessHistory gets access history
func (s *AccessService) GetAccessHistory(ctx context.Context, req GetAccessHistoryRequest) ([]models.AccessLog, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.GetAccessLogs(ctx, req.UserID, req.GateID, perPage, (page-1)*perPage)
}

// LogAccessRequest represents log access request
type LogAccessRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	GateID     string `json:"gate_id" binding:"required"`
	AccessType string `json:"access_type" binding:"required,oneof=entry exit"`
	IsAllowed  bool   `json:"is_allowed"`
	Reason     string `json:"reason,omitempty"`
	QRCodeID   string `json:"qr_code_id,omitempty"`
}

// LogAccess logs an access attempt
func (s *AccessService) LogAccess(ctx context.Context, req LogAccessRequest) (*models.AccessLog, error) {
	log := &models.AccessLog{
		UserID:     req.UserID,
		GateID:     req.GateID,
		AccessType: req.AccessType,
		IsAllowed:  req.IsAllowed,
		Reason:     req.Reason,
	}

	if req.QRCodeID != "" {
		log.QRCodeID = &req.QRCodeID
	}

	if err := s.repo.CreateAccessLog(ctx, log); err != nil {
		return nil, apperrors.NewInternalError("failed to log access", err)
	}

	return log, nil
}

// GetAccessPermissions gets user access permissions
func (s *AccessService) GetAccessPermissions(ctx context.Context, userID string) ([]models.AccessPermission, error) {
	return s.repo.GetUserAccessPermissions(ctx, userID)
}

// CreateAccessPermissionRequest represents create access permission request
type CreateAccessPermissionRequest struct {
	UserID     string  `json:"user_id" binding:"required"`
	GateID     string  `json:"gate_id" binding:"required"`
	IsAllowed  bool    `json:"is_allowed"`
	ValidFrom  *string `json:"valid_from,omitempty"`
	ValidUntil *string `json:"valid_until,omitempty"`
}

// CreateAccessPermission creates an access permission
func (s *AccessService) CreateAccessPermission(ctx context.Context, req CreateAccessPermissionRequest) (*models.AccessPermission, error) {
	permission := &models.AccessPermission{
		UserID:    req.UserID,
		GateID:    req.GateID,
		IsAllowed: req.IsAllowed,
	}

	if req.ValidFrom != nil {
		validFrom, err := time.Parse(time.RFC3339, *req.ValidFrom)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid valid_from format")
		}
		permission.ValidFrom = &validFrom
	}

	if req.ValidUntil != nil {
		validUntil, err := time.Parse(time.RFC3339, *req.ValidUntil)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid valid_until format")
		}
		permission.ValidUntil = &validUntil
	}

	if err := s.repo.CreateAccessPermission(ctx, permission); err != nil {
		return nil, apperrors.NewInternalError("failed to create permission", err)
	}

	return permission, nil
}
