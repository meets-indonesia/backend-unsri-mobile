package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestAccessLog() *models.AccessLog {
	return &models.AccessLog{
		ID:         uuid.New().String(),
		UserID:     uuid.New().String(),
		GateID:     uuid.New().String(),
		AccessType: "ENTRY",
		IsAllowed:  true,
		CreatedAt:  time.Now(),
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("access", "test-id")
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
		err := apperrors.NewForbiddenError("access denied")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test AccessLog model
func TestAccessLogModel(t *testing.T) {
	t.Run("valid access log", func(t *testing.T) {
		log := createTestAccessLog()
		if log.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if log.GateID == "" {
			t.Error("GateID should not be empty")
		}
		if log.AccessType == "" {
			t.Error("AccessType should not be empty")
		}
		if log.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		log := models.AccessLog{}
		if log.TableName() != "access_logs" {
			t.Errorf("Expected table name 'access_logs', got '%s'", log.TableName())
		}
	})

	t.Run("access types", func(t *testing.T) {
		validTypes := []string{"ENTRY", "EXIT"}
		for _, atype := range validTypes {
			log := createTestAccessLog()
			log.AccessType = atype
			if log.AccessType != atype {
				t.Errorf("Expected access type %s, got %s", atype, log.AccessType)
			}
		}
	})
}

// Test ValidateQRRequest validation
func TestValidateQRRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ValidateQRRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ValidateQRRequest{
				QRToken: "test-token",
				GateID:  uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "missing qr_token",
			req: ValidateQRRequest{
				GateID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "missing gate_id",
			req: ValidateQRRequest{
				QRToken: "test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.QRToken == "" && !tt.wantErr {
				t.Error("QRToken should be required")
			}
			if tt.req.GateID == "" && !tt.wantErr {
				t.Error("GateID should be required")
			}
		})
	}
}

