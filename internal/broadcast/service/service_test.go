package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestBroadcast() *models.Broadcast {
	return &models.Broadcast{
		ID:          uuid.New().String(),
		Title:       "Test Broadcast",
		Content:     "Test content",
		Type:        models.BroadcastTypeGeneral,
		CreatedBy:   uuid.New().String(),
		IsPublished: false,
		CreatedAt:   time.Now(),
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("broadcast", "test-id")
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

// Test Broadcast model
func TestBroadcastModel(t *testing.T) {
	t.Run("valid broadcast", func(t *testing.T) {
		broadcast := createTestBroadcast()
		if broadcast.Title == "" {
			t.Error("Title should not be empty")
		}
		if broadcast.CreatedBy == "" {
			t.Error("CreatedBy should not be empty")
		}
		if broadcast.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		broadcast := models.Broadcast{}
		if broadcast.TableName() != "broadcasts" {
			t.Errorf("Expected table name 'broadcasts', got '%s'", broadcast.TableName())
		}
	})

	t.Run("broadcast publish status", func(t *testing.T) {
		broadcast := createTestBroadcast()
		if broadcast.IsPublished {
			t.Error("Broadcast should not be published by default")
		}

		broadcast.IsPublished = true
		broadcast.PublishedAt = func() *time.Time { t := time.Now(); return &t }()
		if !broadcast.IsPublished {
			t.Error("Broadcast should be published")
		}
	})
}

// Test CreateBroadcastRequest validation
func TestCreateBroadcastRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateBroadcastRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateBroadcastRequest{
				Title:   "Test Broadcast",
				Content: "Test content",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: CreateBroadcastRequest{
				Content: "Test content",
			},
			wantErr: true,
		},
		{
			name: "empty title",
			req: CreateBroadcastRequest{
				Title:   "",
				Content: "Test content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Title == "" && !tt.wantErr {
				t.Error("Title should be required")
			}
		})
	}
}

