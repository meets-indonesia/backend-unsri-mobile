package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// AccessRepository handles access control data operations
type AccessRepository struct {
	db *gorm.DB
}

// NewAccessRepository creates a new access repository
func NewAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

// CreateAccessLog creates an access log entry
func (r *AccessRepository) CreateAccessLog(ctx context.Context, log *models.AccessLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetAccessLogs gets access logs with filters
func (r *AccessRepository) GetAccessLogs(ctx context.Context, userID *string, gateID *string, limit, offset int) ([]models.AccessLog, int64, error) {
	var logs []models.AccessLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AccessLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if gateID != nil {
		query = query.Where("gate_id = ?", *gateID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAccessPermission gets access permission for user and gate
func (r *AccessRepository) GetAccessPermission(ctx context.Context, userID, gateID string) (*models.AccessPermission, error) {
	var permission models.AccessPermission
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND gate_id = ? AND is_active = ?", userID, gateID, true).
		Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_until IS NULL OR valid_until >= ?)", now, now).
		First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

// CreateAccessPermission creates an access permission
func (r *AccessRepository) CreateAccessPermission(ctx context.Context, permission *models.AccessPermission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetUserAccessPermissions gets all permissions for a user
func (r *AccessRepository) GetUserAccessPermissions(ctx context.Context, userID string) ([]models.AccessPermission, error) {
	var permissions []models.AccessPermission
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_until IS NULL OR valid_until >= ?)", now, now).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// UpdateAccessPermission updates an access permission
func (r *AccessRepository) UpdateAccessPermission(ctx context.Context, permission *models.AccessPermission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// DeleteAccessPermission soft deletes an access permission
func (r *AccessRepository) DeleteAccessPermission(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AccessPermission{}, "id = ?", id).Error
}

