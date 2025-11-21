package service

import (
	"context"
	"encoding/json"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/attendance/repository"
	"unsri-backend/pkg/jwt"
	"unsri-backend/pkg/qrcode"
)

// AttendanceService handles attendance business logic
type AttendanceService struct {
	repo *repository.AttendanceRepository
	jwt  *jwt.JWT
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService(repo *repository.AttendanceRepository, jwtToken *jwt.JWT) *AttendanceService {
	return &AttendanceService{
		repo: repo,
		jwt:  jwtToken,
	}
}

// GenerateQRRequest represents request to generate QR code
type GenerateQRRequest struct {
	ScheduleID *string `json:"schedule_id,omitempty"`
	Type       string  `json:"type" binding:"required,oneof=kelas kampus"`
	Duration   int     `json:"duration"` // Duration in minutes, default 15
}

// GenerateQRResponse represents QR code generation response
type GenerateQRResponse struct {
	SessionID string `json:"session_id"`
	QRCode    string `json:"qr_code"` // Base64 encoded QR code image
	ExpiresAt string `json:"expires_at"`
}

// GenerateQRCode generates a QR code for attendance
func (s *AttendanceService) GenerateQRCode(ctx context.Context, userID string, req GenerateQRRequest) (*GenerateQRResponse, error) {
	// Validate user role (only dosen/staff can generate QR)
	// This should be validated in middleware, but we check here too

	duration := 15 // Default 15 minutes
	if req.Duration > 0 {
		duration = req.Duration
	}

	expiresAt := time.Now().Add(time.Duration(duration) * time.Minute)

	session := &models.AttendanceSession{
		ScheduleID: req.ScheduleID,
		CreatedBy:  userID,
		Type:       models.AttendanceType(req.Type),
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, apperrors.NewInternalError("failed to create attendance session", err)
	}

	// Generate QR code data
	qrData := qrcode.QRData{
		SessionID:  session.ID,
		ScheduleID: "",
		ExpiresAt:  expiresAt,
		Type:       req.Type,
	}

	if req.ScheduleID != nil {
		qrData.ScheduleID = *req.ScheduleID
	}

	// Generate QR code image
	qrImage, err := qrcode.GenerateQRCode(qrData)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate QR code", err)
	}

	// Store QR code data in session
	qrDataJSON, _ := json.Marshal(qrData)
	session.QRCode = string(qrDataJSON)
	s.repo.UpdateSession(ctx, session) // Ignore error, QR code already generated

	return &GenerateQRResponse{
		SessionID: session.ID,
		QRCode:    string(qrImage), // In production, return base64 or URL
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// ScanQRRequest represents request to scan QR code
type ScanQRRequest struct {
	QRData string  `json:"qr_data" binding:"required"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// ScanQRResponse represents QR scan response
type ScanQRResponse struct {
	AttendanceID string `json:"attendance_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// ScanQRCode scans a QR code and records attendance
func (s *AttendanceService) ScanQRCode(ctx context.Context, userID string, req ScanQRRequest) (*ScanQRResponse, error) {
	// Parse QR data
	qrData, err := qrcode.ParseQRData(req.QRData)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid QR code data")
	}

	// Get session
	session, err := s.repo.GetSessionByID(ctx, qrData.SessionID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid or expired QR code")
	}

	// Check if session is active and not expired
	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, apperrors.NewBadRequestError("QR code has expired")
	}

	// Check if attendance already exists
	date := time.Now()
	exists, err := s.repo.CheckAttendanceExists(ctx, userID, date, session.ScheduleID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check attendance", err)
	}

	if exists {
		return nil, apperrors.NewConflictError("attendance already recorded for today")
	}

	// Create attendance record
	attendance := &models.Attendance{
		UserID:     userID,
		SessionID:  &session.ID,
		ScheduleID: session.ScheduleID,
		Type:       session.Type,
		Status:     models.StatusHadir,
		Date:       date,
		CheckInTime: &date,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
	}

	if err := s.repo.CreateAttendance(ctx, attendance); err != nil {
		return nil, apperrors.NewInternalError("failed to record attendance", err)
	}

	// If this is a class attendance QR, deactivate the session so QR will regenerate
	// The QR service will handle regeneration when generating new QR for the schedule
	if session.Type == models.AttendanceTypeKelas && session.ScheduleID != nil {
		session.IsActive = false
		s.repo.UpdateSession(ctx, session) // Deactivate so new QR can be generated
	}

	return &ScanQRResponse{
		AttendanceID: attendance.ID,
		Status:       string(attendance.Status),
		Message:      "Attendance recorded successfully",
	}, nil
}

// GetAttendancesRequest represents request to get attendances
type GetAttendancesRequest struct {
	UserID    *string `form:"user_id"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
	Page      int     `form:"page,default=1"`
	PerPage   int     `form:"per_page,default=20"`
}

// GetAttendances gets attendance records
func (s *AttendanceService) GetAttendances(ctx context.Context, req GetAttendancesRequest) ([]models.Attendance, int64, error) {
	var startDate, endDate *time.Time
	var userID string

	if req.UserID != nil {
		userID = *req.UserID
	}

	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			startDate = &t
		}
	}

	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &t
		}
	}

	page := req.Page
	if page < 1 {
		page = 1
	}

	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	attendances, total, err := s.repo.GetAttendancesByUserID(ctx, userID, startDate, endDate, perPage, offset)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to get attendances", err)
	}

	return attendances, total, nil
}

// ManualAttendanceRequest represents manual attendance entry request
type ManualAttendanceRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ScheduleID *string `json:"schedule_id,omitempty"`
	Type      string  `json:"type" binding:"required,oneof=kelas kampus"`
	Status    string  `json:"status" binding:"required,oneof=hadir izin sakit alpa terlambat"`
	Date      string  `json:"date" binding:"required"`
	Notes     string  `json:"notes,omitempty"`
}

// CreateManualAttendance creates a manual attendance record
func (s *AttendanceService) CreateManualAttendance(ctx context.Context, createdBy string, req ManualAttendanceRequest) (*models.Attendance, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid date format, use YYYY-MM-DD")
	}

	// Check if attendance already exists
	exists, err := s.repo.CheckAttendanceExists(ctx, req.UserID, date, req.ScheduleID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check attendance", err)
	}

	if exists {
		return nil, apperrors.NewConflictError("attendance already exists for this date")
	}

	attendance := &models.Attendance{
		UserID:     req.UserID,
		ScheduleID: req.ScheduleID,
		Type:       models.AttendanceType(req.Type),
		Status:     models.AttendanceStatus(req.Status),
		Date:       date,
		Notes:      req.Notes,
		CreatedBy:  &createdBy,
	}

	if err := s.repo.CreateAttendance(ctx, attendance); err != nil {
		return nil, apperrors.NewInternalError("failed to create attendance", err)
	}

	return attendance, nil
}

// UpdateAttendanceRequest represents request to update attendance
type UpdateAttendanceRequest struct {
	Status string  `json:"status" binding:"required,oneof=hadir izin sakit alpa terlambat"`
	Notes  string  `json:"notes,omitempty"`
}

// UpdateAttendance updates an attendance record
func (s *AttendanceService) UpdateAttendance(ctx context.Context, attendanceID string, req UpdateAttendanceRequest) (*models.Attendance, error) {
	attendance, err := s.repo.GetAttendanceByID(ctx, attendanceID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("attendance", attendanceID)
	}

	attendance.Status = models.AttendanceStatus(req.Status)
	if req.Notes != "" {
		attendance.Notes = req.Notes
	}

	if err := s.repo.UpdateAttendance(ctx, attendance); err != nil {
		return nil, apperrors.NewInternalError("failed to update attendance", err)
	}

	return attendance, nil
}

// GetAttendanceStatistics gets attendance statistics
func (s *AttendanceService) GetAttendanceStatistics(ctx context.Context, userID string, startDate, endDate *string) (map[string]interface{}, error) {
	var start, end *time.Time

	if startDate != nil {
		if t, err := time.Parse("2006-01-02", *startDate); err == nil {
			start = &t
		}
	}

	if endDate != nil {
		if t, err := time.Parse("2006-01-02", *endDate); err == nil {
			end = &t
		}
	}

	return s.repo.GetAttendanceStatistics(ctx, userID, start, end)
}

// GetAttendanceByCourse gets attendances by course
func (s *AttendanceService) GetAttendanceByCourse(ctx context.Context, courseID string, startDate, endDate *string) ([]models.Attendance, error) {
	var start, end *time.Time

	if startDate != nil {
		if t, err := time.Parse("2006-01-02", *startDate); err == nil {
			start = &t
		}
	}

	if endDate != nil {
		if t, err := time.Parse("2006-01-02", *endDate); err == nil {
			end = &t
		}
	}

	return s.repo.GetAttendancesByCourseID(ctx, courseID, start, end)
}

// GetAttendanceByStudent gets attendances by student
func (s *AttendanceService) GetAttendanceByStudent(ctx context.Context, studentID string, startDate, endDate *string) ([]models.Attendance, error) {
	var start, end *time.Time

	if startDate != nil {
		if t, err := time.Parse("2006-01-02", *startDate); err == nil {
			start = &t
		}
	}

	if endDate != nil {
		if t, err := time.Parse("2006-01-02", *endDate); err == nil {
			end = &t
		}
	}

	return s.repo.GetAttendancesByStudentID(ctx, studentID, start, end)
}

// GetAttendanceOverview gets attendance overview
func (s *AttendanceService) GetAttendanceOverview(ctx context.Context, userID string, role string) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// Today's schedules
	todaySchedules, err := s.repo.GetTodaySchedules(ctx, userID, role)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get today's schedules", err)
	}
	overview["today_schedules"] = todaySchedules

	// Upcoming schedules
	upcomingSchedules, err := s.repo.GetUpcomingSchedules(ctx, userID, role, 5)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get upcoming schedules", err)
	}
	overview["upcoming_schedules"] = upcomingSchedules

	// Statistics for this month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	
	startStr := startOfMonth.Format("2006-01-02")
	endStr := endOfMonth.Format("2006-01-02")
	
	stats, err := s.GetAttendanceStatistics(ctx, userID, &startStr, &endStr)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get statistics", err)
	}
	overview["monthly_statistics"] = stats

	// Current tap in status (for campus attendance)
	if role == "mahasiswa" || role == "staff" {
		tapInStatus, err := s.repo.GetCurrentTapInStatus(ctx, userID)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get tap in status", err)
		}
		overview["current_tap_in"] = tapInStatus
	}

	return overview, nil
}

// TapInRequest represents tap in request
type TapInRequest struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// TapIn performs tap in for campus attendance
func (s *AttendanceService) TapIn(ctx context.Context, userID string, req TapInRequest) (*models.Attendance, error) {
	// Check if already tapped in today
	current, err := s.repo.GetCurrentTapInStatus(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check tap in status", err)
	}

	if current != nil {
		return nil, apperrors.NewConflictError("already tapped in today")
	}

	now := time.Now()
	attendance := &models.Attendance{
		UserID:      userID,
		Type:        models.AttendanceTypeKampus,
		Status:      models.StatusHadir,
		Date:        now,
		CheckInTime: &now,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := s.repo.CreateAttendance(ctx, attendance); err != nil {
		return nil, apperrors.NewInternalError("failed to record tap in", err)
	}

	return attendance, nil
}

// TapOut performs tap out for campus attendance
func (s *AttendanceService) TapOut(ctx context.Context, userID string, req TapInRequest) (*models.Attendance, error) {
	// Get current tap in
	current, err := s.repo.GetCurrentTapInStatus(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to check tap in status", err)
	}

	if current == nil {
		return nil, apperrors.NewBadRequestError("no active tap in found")
	}

	now := time.Now()
	current.CheckOutTime = &now
	if req.Latitude != nil {
		current.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		current.Longitude = req.Longitude
	}

	if err := s.repo.UpdateAttendance(ctx, current); err != nil {
		return nil, apperrors.NewInternalError("failed to record tap out", err)
	}

	return current, nil
}

// CreateScheduleRequest represents request to create schedule
type CreateScheduleRequest struct {
	CourseID   *string `json:"course_id,omitempty"`
	CourseCode string  `json:"course_code"`
	CourseName string  `json:"course_name"`
	DosenID    string  `json:"dosen_id" binding:"required"`
	Room       string  `json:"room"`
	DayOfWeek  int     `json:"day_of_week" binding:"required,min=0,max=6"`
	StartTime  string  `json:"start_time" binding:"required"`
	EndTime    string  `json:"end_time" binding:"required"`
	Date       string  `json:"date" binding:"required"`
}

// CreateSchedule creates a new schedule
func (s *AttendanceService) CreateSchedule(ctx context.Context, req CreateScheduleRequest) (*models.Schedule, error) {
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_time format, use HH:MM")
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_time format, use HH:MM")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid date format, use YYYY-MM-DD")
	}

	// Combine date with time
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, date.Location())
	endDateTime := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, date.Location())

	schedule := &models.Schedule{
		CourseID:   req.CourseID,
		CourseCode: req.CourseCode,
		CourseName: req.CourseName,
		DosenID:    req.DosenID,
		Room:       req.Room,
		DayOfWeek:  req.DayOfWeek,
		StartTime:  startDateTime,
		EndTime:    endDateTime,
		Date:       date,
		IsActive:   true,
	}

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, apperrors.NewInternalError("failed to create schedule", err)
	}

	return schedule, nil
}

// GetSchedulesRequest represents request to get schedules
type GetSchedulesRequest struct {
	DosenID   *string `form:"dosen_id"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
	Page      int     `form:"page,default=1"`
	PerPage   int     `form:"per_page,default=20"`
}

// GetSchedules gets schedules
func (s *AttendanceService) GetSchedules(ctx context.Context, req GetSchedulesRequest) ([]models.Schedule, int64, error) {
	var startDate, endDate *time.Time

	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			startDate = &t
		}
	}

	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &t
		}
	}

	page := req.Page
	if page < 1 {
		page = 1
	}

	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	return s.repo.GetAllSchedules(ctx, req.DosenID, startDate, endDate, perPage, offset)
}

// UpdateScheduleRequest represents request to update schedule
type UpdateScheduleRequest struct {
	CourseCode *string `json:"course_code,omitempty"`
	CourseName *string `json:"course_name,omitempty"`
	Room       *string `json:"room,omitempty"`
	StartTime  *string `json:"start_time,omitempty"`
	EndTime    *string `json:"end_time,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// UpdateSchedule updates a schedule
func (s *AttendanceService) UpdateSchedule(ctx context.Context, scheduleID string, req UpdateScheduleRequest) (*models.Schedule, error) {
	schedule, err := s.repo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("schedule", scheduleID)
	}

	if req.CourseCode != nil {
		schedule.CourseCode = *req.CourseCode
	}
	if req.CourseName != nil {
		schedule.CourseName = *req.CourseName
	}
	if req.Room != nil {
		schedule.Room = *req.Room
	}
	if req.StartTime != nil {
		startTime, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid start_time format")
		}
		schedule.StartTime = time.Date(schedule.Date.Year(), schedule.Date.Month(), schedule.Date.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, schedule.Date.Location())
	}
	if req.EndTime != nil {
		endTime, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid end_time format")
		}
		schedule.EndTime = time.Date(schedule.Date.Year(), schedule.Date.Month(), schedule.Date.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, schedule.Date.Location())
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, apperrors.NewInternalError("failed to update schedule", err)
	}

	return schedule, nil
}

// GetScheduleByID gets a schedule by ID
func (s *AttendanceService) GetScheduleByID(ctx context.Context, scheduleID string) (*models.Schedule, error) {
	return s.repo.GetScheduleByID(ctx, scheduleID)
}

// DeleteSchedule deletes a schedule
func (s *AttendanceService) DeleteSchedule(ctx context.Context, scheduleID string) error {
	_, err := s.repo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return apperrors.NewNotFoundError("schedule", scheduleID)
	}

	return s.repo.DeleteSchedule(ctx, scheduleID)
}

