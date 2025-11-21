package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// NotificationRepository handles notification data operations
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetNotificationByID gets a notification by ID
func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id string) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

// GetNotificationsByUserID gets notifications for a user
func (r *NotificationRepository) GetNotificationsByUserID(ctx context.Context, userID string, isRead *bool, limit, offset int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ?", userID)

	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// UpdateNotification updates a notification
func (r *NotificationRepository) UpdateNotification(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// CreateDeviceToken creates a device token
func (r *NotificationRepository) CreateDeviceToken(ctx context.Context, token *models.DeviceToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetDeviceTokenByToken gets a device token by token string
func (r *NotificationRepository) GetDeviceTokenByToken(ctx context.Context, tokenStr string) (*models.DeviceToken, error) {
	var token models.DeviceToken
	if err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("device token not found")
		}
		return nil, err
	}
	return &token, nil
}

// GetDeviceTokensByUserID gets device tokens for a user
func (r *NotificationRepository) GetDeviceTokensByUserID(ctx context.Context, userID string) ([]models.DeviceToken, error) {
	var tokens []models.DeviceToken
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// DeleteDeviceToken deletes a device token
func (r *NotificationRepository) DeleteDeviceToken(ctx context.Context, tokenStr string) error {
	return r.db.WithContext(ctx).Where("token = ?", tokenStr).Delete(&models.DeviceToken{}).Error
}

// UpdateDeviceToken updates a device token
func (r *NotificationRepository) UpdateDeviceToken(ctx context.Context, token *models.DeviceToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

