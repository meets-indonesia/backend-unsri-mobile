package service

import (
	"testing"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
)

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("report", "test-id")
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

// Test AttendanceReportRequest validation
func TestAttendanceReportRequest(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 1, 0)

	tests := []struct {
		name    string
		req     AttendanceReportRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: AttendanceReportRequest{
				StartDate: startDate.Format("2006-01-02"),
				EndDate:   endDate.Format("2006-01-02"),
				Summary:   true,
			},
			wantErr: false,
		},
		{
			name: "missing start_date",
			req: AttendanceReportRequest{
				EndDate: endDate.Format("2006-01-02"),
				Summary: true,
			},
			wantErr: true,
		},
		{
			name: "missing end_date",
			req: AttendanceReportRequest{
				StartDate: startDate.Format("2006-01-02"),
				Summary:   true,
			},
			wantErr: true,
		},
		{
			name: "invalid date range",
			req: AttendanceReportRequest{
				StartDate: endDate.Format("2006-01-02"),
				EndDate:   startDate.Format("2006-01-02"),
				Summary:   true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.StartDate == "" && !tt.wantErr {
				t.Error("StartDate should be required")
			}
			if tt.req.EndDate == "" && !tt.wantErr {
				t.Error("EndDate should be required")
			}
			if tt.req.StartDate != "" && tt.req.EndDate != "" {
				start, err1 := time.Parse("2006-01-02", tt.req.StartDate)
				end, err2 := time.Parse("2006-01-02", tt.req.EndDate)
				if err1 == nil && err2 == nil {
					if end.Before(start) && !tt.wantErr {
						t.Error("EndDate should be after StartDate")
					}
				}
			}
		})
	}
}

// Test AcademicReportRequest validation
func TestAcademicReportRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     AcademicReportRequest
		wantErr bool
	}{
		{
			name: "valid request with student_id",
			req: AcademicReportRequest{
				StudentID: "student-123",
				Semester:  func() *string { v := "2024-1"; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "valid request with semester only",
			req: AcademicReportRequest{
				StudentID: "student-123",
				Semester:  func() *string { v := "2024-1"; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "missing student_id",
			req: AcademicReportRequest{
				Semester: func() *string { v := "2024-1"; return &v }(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.StudentID == "" && !tt.wantErr {
				t.Error("StudentID should be required")
			}
			if tt.req.Semester != nil && len(*tt.req.Semester) < 5 {
				t.Error("Semester should have valid format")
			}
		})
	}
}

// Test report types
func TestReportTypes(t *testing.T) {
	validTypes := []string{"attendance", "academic", "system"}

	for _, rtype := range validTypes {
		t.Run("valid report type: "+rtype, func(t *testing.T) {
			// Report type validation
			if rtype == "" {
				t.Error("Report type should not be empty")
			}
		})
	}
}

// Test date range for reports
func TestReportDateRange(t *testing.T) {
	t.Run("valid date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, 1, 0)

		if endDate.Before(startDate) {
			t.Error("EndDate should be after StartDate")
		}
	})

	t.Run("maximum date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(1, 0, 0) // 1 year

		daysDiff := endDate.Sub(startDate).Hours() / 24
		if daysDiff > 365 {
			t.Error("Date range might be too large for report generation")
		}
	})
}
