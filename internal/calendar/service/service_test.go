package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestCalendarEvent() *models.AcademicEvent {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)
	return &models.AcademicEvent{
		ID:          uuid.New().String(),
		Title:       "Test Event",
		Description: "Test description",
		StartDate:   startTime,
		EndDate:     endTime,
		CreatedBy:   uuid.New().String(),
		IsActive:    true,
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("calendar", "test-id")
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
		err := apperrors.NewConflictError("event conflict")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test CalendarEvent model
func TestCalendarEventModel(t *testing.T) {
	t.Run("valid calendar event", func(t *testing.T) {
		event := createTestCalendarEvent()
		if event.Title == "" {
			t.Error("Title should not be empty")
		}
		if event.CreatedBy == "" {
			t.Error("CreatedBy should not be empty")
		}
		if event.StartDate.After(event.EndDate) {
			t.Error("StartDate should be before EndDate")
		}
		if event.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		event := models.AcademicEvent{}
		if event.TableName() != "academic_events" {
			t.Errorf("Expected table name 'academic_events', got '%s'", event.TableName())
		}
	})

	t.Run("time range validation", func(t *testing.T) {
		event := createTestCalendarEvent()
		if event.EndDate.Before(event.StartDate) {
			t.Error("EndDate should be after StartDate")
		}
	})
}

// Test CreateEventRequest validation
func TestCreateEventRequest(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	tests := []struct {
		name    string
		req     CreateEventRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartDate: startTime.Format(time.RFC3339),
				EndDate:   endTime.Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: CreateEventRequest{
				StartDate: startTime.Format(time.RFC3339),
				EndDate:   endTime.Format(time.RFC3339),
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartDate: "invalid-date",
				EndDate:   endTime.Format(time.RFC3339),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.req.Title == "" {
				hasError = true
			}
			if tt.req.StartDate != "" {
				_, err := time.Parse(time.RFC3339, tt.req.StartDate)
				if err != nil {
					hasError = true
				}
			}
			if tt.req.EndDate != "" {
				_, err := time.Parse(time.RFC3339, tt.req.EndDate)
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

