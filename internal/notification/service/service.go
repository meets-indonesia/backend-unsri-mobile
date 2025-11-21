package service

import (
	"context"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/notification/repository"
)

// NotificationService handles notification business logic
type NotificationService struct {
	repo *repository.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// SendNotificationRequest represents send notification request
type SendNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=info warning error success"`
	Data    string `json:"data,omitempty"`
}

// SendNotification sends a notification
func (s *NotificationService) SendNotification(ctx context.Context, req SendNotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    models.NotificationType(req.Type),
		Data:    req.Data,
		IsRead:  false,
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return nil, apperrors.NewInternalError("failed to create notification", err)
	}

	// TODO: Send push notification via FCM

	return notification, nil
}

// GetNotificationsRequest represents get notifications request
type GetNotificationsRequest struct {
	IsRead *bool `form:"is_read"`
	Page   int   `form:"page,default=1"`
	PerPage int  `form:"per_page,default=20"`
}

// GetNotifications gets notifications for a user
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, req GetNotificationsRequest) ([]models.Notification, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.GetNotificationsByUserID(ctx, userID, req.IsRead, perPage, (page-1)*perPage)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	notification, err := s.repo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return apperrors.NewNotFoundError("notification", notificationID)
	}

	if notification.UserID != userID {
		return apperrors.NewForbiddenError("not authorized to read this notification")
	}

	if !notification.IsRead {
		now := time.Now()
		notification.IsRead = true
		notification.ReadAt = &now
		if err := s.repo.UpdateNotification(ctx, notification); err != nil {
			return apperrors.NewInternalError("failed to update notification", err)
		}
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// RegisterDeviceTokenRequest represents register device token request
type RegisterDeviceTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web"`
}

// RegisterDeviceToken registers a device token for push notifications
func (s *NotificationService) RegisterDeviceToken(ctx context.Context, userID string, req RegisterDeviceTokenRequest) (*models.DeviceToken, error) {
	// Check if token already exists
	existing, _ := s.repo.GetDeviceTokenByToken(ctx, req.Token)
	if existing != nil {
		// Update if exists
		existing.UserID = userID
		existing.Platform = req.Platform
		existing.IsActive = true
		if err := s.repo.UpdateDeviceToken(ctx, existing); err != nil {
			return nil, apperrors.NewInternalError("failed to update device token", err)
		}
		return existing, nil
	}

	token := &models.DeviceToken{
		UserID:   userID,
		Token:    req.Token,
		Platform: req.Platform,
		IsActive: true,
	}

	if err := s.repo.CreateDeviceToken(ctx, token); err != nil {
		return nil, apperrors.NewInternalError("failed to register device token", err)
	}

	return token, nil
}

// UnregisterDeviceToken unregisters a device token
func (s *NotificationService) UnregisterDeviceToken(ctx context.Context, tokenStr string) error {
	return s.repo.DeleteDeviceToken(ctx, tokenStr)
}

// GetDeviceTokens gets device tokens for a user
func (s *NotificationService) GetDeviceTokens(ctx context.Context, userID string) ([]models.DeviceToken, error) {
	return s.repo.GetDeviceTokensByUserID(ctx, userID)
}

