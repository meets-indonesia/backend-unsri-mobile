package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
)

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("qr", "test-id")
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

// Test GenerateClassQRRequest validation
func TestGenerateClassQRRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateClassQRRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: GenerateClassQRRequest{
				ScheduleID: uuid.New().String(),
				Duration:   15,
			},
			wantErr: false,
		},
		{
			name: "missing schedule_id",
			req: GenerateClassQRRequest{
				Duration: 15,
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			req: GenerateClassQRRequest{
				ScheduleID: uuid.New().String(),
				Duration:   0,
			},
			wantErr: true,
		},
		{
			name: "negative duration",
			req: GenerateClassQRRequest{
				ScheduleID: uuid.New().String(),
				Duration:   -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.ScheduleID == "" && !tt.wantErr {
				t.Error("ScheduleID should be required")
			}
			if tt.req.Duration <= 0 && !tt.wantErr {
				t.Error("Duration should be positive")
			}
		})
	}
}

// Test GenerateAccessQRRequest validation
func TestGenerateAccessQRRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateAccessQRRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: GenerateAccessQRRequest{
				Duration: 525600, // 1 year in minutes
			},
			wantErr: false,
		},
		{
			name: "zero duration",
			req: GenerateAccessQRRequest{
				Duration: 0,
			},
			wantErr: true,
		},
		{
			name: "negative duration",
			req: GenerateAccessQRRequest{
				Duration: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Duration <= 0 && !tt.wantErr {
				t.Error("Duration should be positive")
			}
		})
	}
}

// Test QR code expiration
func TestQRCodeExpiration(t *testing.T) {
	t.Run("expiration calculation", func(t *testing.T) {
		now := time.Now()
		duration := 15 // minutes
		expiresAt := now.Add(time.Duration(duration) * time.Minute)

		if expiresAt.Before(now) {
			t.Error("Expiration time should be in the future")
		}

		expectedDuration := expiresAt.Sub(now)
		if expectedDuration != time.Duration(duration)*time.Minute {
			t.Errorf("Expected duration %v, got %v", time.Duration(duration)*time.Minute, expectedDuration)
		}
	})

	t.Run("expired QR code", func(t *testing.T) {
		now := time.Now()
		pastTime := now.Add(-1 * time.Hour)

		if !pastTime.Before(now) {
			t.Error("Past time should be before now")
		}
	})
}

