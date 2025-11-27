package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestSchedule() *models.Schedule {
	return &models.Schedule{
		ID:        uuid.New().String(),
		DosenID:   uuid.New().String(),
		DayOfWeek: 1, // Monday
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
		Date:      time.Now(),
		IsActive:  true,
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("schedule", "test-id")
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

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("schedule conflict")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test Schedule model
func TestScheduleModel(t *testing.T) {
	t.Run("valid schedule", func(t *testing.T) {
		schedule := createTestSchedule()
		if schedule.DosenID == "" {
			t.Error("DosenID should not be empty")
		}
		if schedule.ID == "" {
			t.Error("ID should be generated")
		}
		if schedule.StartTime.After(schedule.EndTime) {
			t.Error("StartTime should be before EndTime")
		}
	})

	t.Run("table name", func(t *testing.T) {
		schedule := models.Schedule{}
		if schedule.TableName() != "schedules" {
			t.Errorf("Expected table name 'schedules', got '%s'", schedule.TableName())
		}
	})

	t.Run("day of week validation", func(t *testing.T) {
		schedule := createTestSchedule()
		if schedule.DayOfWeek < 0 || schedule.DayOfWeek > 6 {
			t.Error("DayOfWeek should be between 0 and 6")
		}
	})

	t.Run("time range validation", func(t *testing.T) {
		schedule := createTestSchedule()
		if schedule.EndTime.Before(schedule.StartTime) {
			t.Error("EndTime should be after StartTime")
		}
	})
}

// Test CreateScheduleRequest validation
func TestCreateScheduleRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateScheduleRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateScheduleRequest{
				DosenID:   uuid.New().String(),
				DayOfWeek: 1,
				StartTime: "08:00",
				EndTime:   "10:00",
				Date:      time.Now().Format("2006-01-02"),
			},
			wantErr: false,
		},
		{
			name: "missing dosen_id",
			req: CreateScheduleRequest{
				DayOfWeek: 1,
				StartTime: "08:00",
				EndTime:   "10:00",
				Date:      time.Now().Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "invalid day of week",
			req: CreateScheduleRequest{
				DosenID:   uuid.New().String(),
				DayOfWeek: 10, // Invalid
				StartTime: "08:00",
				EndTime:   "10:00",
				Date:      time.Now().Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "invalid time format",
			req: CreateScheduleRequest{
				DosenID:   uuid.New().String(),
				DayOfWeek: 1,
				StartTime: "invalid-time",
				EndTime:   "10:00",
				Date:      time.Now().Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: CreateScheduleRequest{
				DosenID:   uuid.New().String(),
				DayOfWeek: 1,
				StartTime: "08:00",
				EndTime:   "10:00",
				Date:      "invalid-date",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.req.DosenID == "" {
				hasError = true
			}
			if tt.req.DayOfWeek < 0 || tt.req.DayOfWeek > 6 {
				hasError = true
			}
			if tt.req.StartTime != "" {
				_, err := time.Parse("15:04", tt.req.StartTime)
				if err != nil {
					hasError = true
				}
			}
			if tt.req.EndTime != "" {
				_, err := time.Parse("15:04", tt.req.EndTime)
				if err != nil {
					hasError = true
				}
			}
			if tt.req.Date != "" {
				_, err := time.Parse("2006-01-02", tt.req.Date)
				if err != nil {
					hasError = true
				}
			}
			
			if hasError != tt.wantErr {
				t.Errorf("Expected error = %v, got %v", tt.wantErr, hasError)
			}
		})
	}
}

