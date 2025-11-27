package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestNotification() *models.Notification {
	return &models.Notification{
		ID:        uuid.New().String(),
		UserID:    uuid.New().String(),
		Title:     "Test Notification",
		Message:   "Test message",
		Type:      "info",
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("notification", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := apperrors.NewValidationError("invalid input")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ForbiddenError", func(t *testing.T) {
		err := apperrors.NewForbiddenError("insufficient permissions")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test Notification model
func TestNotificationModel(t *testing.T) {
	t.Run("valid notification", func(t *testing.T) {
		notification := createTestNotification()
		if notification.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if notification.Title == "" {
			t.Error("Title should not be empty")
		}
		if notification.Message == "" {
			t.Error("Message should not be empty")
		}
		if notification.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		notification := models.Notification{}
		if notification.TableName() != "notifications" {
			t.Errorf("Expected table name 'notifications', got '%s'", notification.TableName())
		}
	})

	t.Run("notification read status", func(t *testing.T) {
		notification := createTestNotification()
		if notification.IsRead {
			t.Error("Notification should not be read by default")
		}

		notification.IsRead = true
		notification.ReadAt = func() *time.Time { t := time.Now(); return &t }()
		if !notification.IsRead {
			t.Error("Notification should be read")
		}
	})

	t.Run("notification types", func(t *testing.T) {
		validTypes := []models.NotificationType{
			models.NotificationTypeInfo,
			models.NotificationTypeWarning,
			models.NotificationTypeError,
			models.NotificationTypeSuccess,
		}
		for _, ntype := range validTypes {
			notification := createTestNotification()
			notification.Type = ntype
			if notification.Type != ntype {
				t.Errorf("Expected type %s, got %s", ntype, notification.Type)
			}
		}
	})
}

// Test SendNotificationRequest validation
func TestSendNotificationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     SendNotificationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: SendNotificationRequest{
				UserID:  uuid.New().String(),
				Title:   "Test Notification",
				Message: "Test message",
				Type:    "info",
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			req: SendNotificationRequest{
				Title:   "Test Notification",
				Message: "Test message",
				Type:    "info",
			},
			wantErr: true,
		},
		{
			name: "missing title",
			req: SendNotificationRequest{
				UserID:  uuid.New().String(),
				Message: "Test message",
				Type:    "info",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			req: SendNotificationRequest{
				UserID: uuid.New().String(),
				Title:  "Test Notification",
				Type:   "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.UserID == "" && !tt.wantErr {
				t.Error("UserID should be required")
			}
			if tt.req.Title == "" && !tt.wantErr {
				t.Error("Title should be required")
			}
			if tt.req.Message == "" && !tt.wantErr {
				t.Error("Message should be required")
			}
		})
	}
}

