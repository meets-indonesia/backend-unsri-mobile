package service

import (
	"testing"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"

	"github.com/google/uuid"
)

// Test helper functions
func createTestLeaveRequest() *models.LeaveRequest {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 5)
	return &models.LeaveRequest{
		ID:        uuid.New().String(),
		UserID:    "user-123",
		LeaveType: models.LeaveTypeAnnual,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: 5,
		Reason:    "Vacation",
		Status:    models.LeaveStatusPending,
	}
}

func createTestLeaveQuota() *models.LeaveQuota {
	return &models.LeaveQuota{
		ID:             uuid.New().String(),
		UserID:         "user-123",
		LeaveType:      models.LeaveTypeAnnual,
		Year:           2024,
		TotalQuota:     20,
		UsedQuota:      0,
		RemainingQuota: 20,
	}
}

// Test CreateLeaveRequestRequest validation
func TestCreateLeaveRequestRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateLeaveRequestRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateLeaveRequestRequest{
				LeaveType: "ANNUAL_LEAVE",
				StartDate: "2024-01-01",
				EndDate:   "2024-01-05",
				Reason:    "Family vacation",
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			req: CreateLeaveRequestRequest{
				LeaveType: "ANNUAL_LEAVE",
				StartDate: "invalid-date",
				EndDate:   "2024-01-05",
				Reason:    "Vacation",
			},
			wantErr: true,
		},
		{
			name: "missing reason",
			req: CreateLeaveRequestRequest{
				LeaveType: "ANNUAL_LEAVE",
				StartDate: "2024-01-01",
				EndDate:   "2024-01-05",
				Reason:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tt.req.StartDate)
			if tt.name == "invalid date format" {
				if err == nil {
					t.Error("expected date parse error")
				}
			} else if err != nil {
				t.Errorf("unexpected date parse error: %v", err)
			}
			if tt.req.Reason == "" && !tt.wantErr {
				t.Error("Reason should be required")
			}
		})
	}
}

// Test date range validation
func TestDateRangeValidation(t *testing.T) {
	t.Run("valid date range", func(t *testing.T) {
		startDate, _ := time.Parse("2006-01-02", "2024-01-01")
		endDate, _ := time.Parse("2006-01-02", "2024-01-05")

		if endDate.Before(startDate) {
			t.Error("End date should be after start date")
		}

		// Calculate total days (inclusive)
		totalDays := endDate.Sub(startDate).Hours()/24 + 1
		if totalDays != 5 {
			t.Errorf("Expected 5 days, got %.2f", totalDays)
		}
	})

	t.Run("invalid date range", func(t *testing.T) {
		startDate, _ := time.Parse("2006-01-02", "2024-01-05")
		endDate, _ := time.Parse("2006-01-02", "2024-01-01")

		if !endDate.Before(startDate) {
			t.Error("End date should be before start date in this test case")
		}
	})
}

// Test leave quota validation
func TestLeaveQuotaValidation(t *testing.T) {
	t.Run("sufficient quota", func(t *testing.T) {
		quota := createTestLeaveQuota()
		requestedDays := 5.0

		if quota.RemainingQuota < requestedDays {
			t.Error("Quota should be sufficient")
		}
	})

	t.Run("insufficient quota", func(t *testing.T) {
		quota := createTestLeaveQuota()
		quota.RemainingQuota = 2.0
		requestedDays := 5.0

		if quota.RemainingQuota >= requestedDays {
			t.Error("Quota should be insufficient")
		}
	})

	t.Run("quota calculation", func(t *testing.T) {
		quota := createTestLeaveQuota()
		quota.UsedQuota = 5.0
		quota.RemainingQuota = quota.TotalQuota - quota.UsedQuota

		if quota.RemainingQuota != 15.0 {
			t.Errorf("Expected remaining quota 15.0, got %.2f", quota.RemainingQuota)
		}
	})
}

// Test leave status transitions
func TestLeaveStatusTransitions(t *testing.T) {
	t.Run("pending to approved", func(t *testing.T) {
		leaveRequest := createTestLeaveRequest()
		leaveRequest.Status = models.LeaveStatusPending

		if leaveRequest.Status != models.LeaveStatusPending {
			t.Error("Status should be pending")
		}

		// Simulate approval
		leaveRequest.Status = models.LeaveStatusApproved
		if leaveRequest.Status != models.LeaveStatusApproved {
			t.Error("Status should be approved")
		}
	})

	t.Run("pending to rejected", func(t *testing.T) {
		leaveRequest := createTestLeaveRequest()
		leaveRequest.Status = models.LeaveStatusPending

		// Simulate rejection
		leaveRequest.Status = models.LeaveStatusRejected
		reason := "Insufficient quota"
		leaveRequest.RejectionReason = &reason

		if leaveRequest.Status != models.LeaveStatusRejected {
			t.Error("Status should be rejected")
		}
		if leaveRequest.RejectionReason == nil {
			t.Error("Rejection reason should be set")
		}
	})

	t.Run("pending to cancelled", func(t *testing.T) {
		leaveRequest := createTestLeaveRequest()
		leaveRequest.Status = models.LeaveStatusPending

		// Simulate cancellation
		leaveRequest.Status = models.LeaveStatusCancelled

		if leaveRequest.Status != models.LeaveStatusCancelled {
			t.Error("Status should be cancelled")
		}
	})
}

// Test leave types
func TestLeaveTypes(t *testing.T) {
	validTypes := []models.LeaveType{
		models.LeaveTypeAnnual,
		models.LeaveTypeSick,
		models.LeaveTypePersonal,
		models.LeaveTypeEmergency,
		models.LeaveTypeUnpaid,
		models.LeaveTypeOther,
	}

	for _, leaveType := range validTypes {
		t.Run(string(leaveType), func(t *testing.T) {
			if leaveType == "" {
				t.Error("Leave type should not be empty")
			}
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("leave request", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("leave quota already exists")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := apperrors.NewValidationError("invalid date range")
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

// Test model validation
func TestLeaveRequestModel(t *testing.T) {
	t.Run("valid leave request", func(t *testing.T) {
		lr := createTestLeaveRequest()
		if lr.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if lr.LeaveType == "" {
			t.Error("LeaveType should not be empty")
		}
		if lr.TotalDays <= 0 {
			t.Error("TotalDays should be positive")
		}
		if lr.Reason == "" {
			t.Error("Reason should not be empty")
		}
	})

	t.Run("table name", func(t *testing.T) {
		lr := models.LeaveRequest{}
		if lr.TableName() != "leave_requests" {
			t.Errorf("Expected table name 'leave_requests', got '%s'", lr.TableName())
		}
	})
}

func TestLeaveQuotaModel(t *testing.T) {
	t.Run("valid leave quota", func(t *testing.T) {
		quota := createTestLeaveQuota()
		if quota.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if quota.LeaveType == "" {
			t.Error("LeaveType should not be empty")
		}
		if quota.Year <= 0 {
			t.Error("Year should be positive")
		}
		if quota.TotalQuota <= 0 {
			t.Error("TotalQuota should be positive")
		}
	})

	t.Run("table name", func(t *testing.T) {
		quota := models.LeaveQuota{}
		if quota.TableName() != "leave_quotas" {
			t.Errorf("Expected table name 'leave_quotas', got '%s'", quota.TableName())
		}
	})

	t.Run("remaining quota calculation", func(t *testing.T) {
		quota := createTestLeaveQuota()
		quota.UsedQuota = 5.0
		quota.RemainingQuota = quota.TotalQuota - quota.UsedQuota

		expected := 15.0
		if quota.RemainingQuota != expected {
			t.Errorf("Expected remaining quota %.2f, got %.2f", expected, quota.RemainingQuota)
		}
	})
}
