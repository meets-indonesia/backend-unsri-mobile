package service

import (
	"context"
	"time"

	"unsri-backend/internal/report/repository"
	apperrors "unsri-backend/internal/shared/errors"
)

// ReportService handles report business logic
type ReportService struct {
	repo *repository.ReportRepository
}

// NewReportService creates a new report service
func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// AttendanceReportRequest represents attendance report request
type AttendanceReportRequest struct {
	StudentID *string `form:"student_id"`
	CourseID  *string `form:"course_id"`
	StartDate string  `form:"start_date" binding:"required"`
	EndDate   string  `form:"end_date" binding:"required"`
	Summary   bool    `form:"summary,default=false"`
}

// AttendanceReportResponse represents attendance report response
type AttendanceReportResponse struct {
	Type    string      `json:"type"`
	Period  Period      `json:"period"`
	Data    interface{} `json:"data"`
	Summary interface{} `json:"summary,omitempty"`
}

// Period represents time period
type Period struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// GetAttendanceReport gets attendance report
func (s *ReportService) GetAttendanceReport(ctx context.Context, req AttendanceReportRequest) (*AttendanceReportResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_date format, use YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return nil, apperrors.NewValidationError("end_date must be after start_date")
	}

	attendances, err := s.repo.GetAttendanceReport(ctx, req.StudentID, req.CourseID, startDate, endDate)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get attendance report", err)
	}

	response := &AttendanceReportResponse{
		Type: "attendance",
		Period: Period{
			StartDate: startDate,
			EndDate:   endDate,
		},
		Data: attendances,
	}

	if req.Summary {
		summary, err := s.repo.GetAttendanceSummary(ctx, req.StudentID, req.CourseID, startDate, endDate)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get attendance summary", err)
		}
		response.Summary = summary
	}

	return response, nil
}

// AcademicReportRequest represents academic report request
type AcademicReportRequest struct {
	StudentID string  `form:"student_id" binding:"required"`
	Semester  *string `form:"semester"`
}

// GetAcademicReport gets academic report for a student
func (s *ReportService) GetAcademicReport(ctx context.Context, req AcademicReportRequest) (map[string]interface{}, error) {
	return s.repo.GetStudentAcademicReport(ctx, req.StudentID, req.Semester)
}

// CourseReportRequest represents course report request
type CourseReportRequest struct {
	CourseID  string `form:"course_id" binding:"required"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// GetCourseReport gets report for a course
func (s *ReportService) GetCourseReport(ctx context.Context, req CourseReportRequest) (map[string]interface{}, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_date format, use YYYY-MM-DD")
	}

	return s.repo.GetCourseReport(ctx, req.CourseID, startDate, endDate)
}

// DailyReportRequest represents daily report request
type DailyReportRequest struct {
	Date string `form:"date"` // YYYY-MM-DD, defaults to today
}

// GetDailyReport gets daily report
func (s *ReportService) GetDailyReport(ctx context.Context, req DailyReportRequest) (map[string]interface{}, error) {
	var date time.Time
	var err error

	if req.Date != "" {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid date format, use YYYY-MM-DD")
		}
	} else {
		date = time.Now()
	}

	return s.repo.GetDailyReport(ctx, date)
}
