package service

import (
	"context"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/schedule/repository"
)

// ScheduleService handles schedule business logic
type ScheduleService struct {
	repo *repository.ScheduleRepository
}

// NewScheduleService creates a new schedule service
func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

// CreateScheduleRequest represents create schedule request
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
func (s *ScheduleService) CreateSchedule(ctx context.Context, req CreateScheduleRequest) (*models.Schedule, error) {
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

// GetScheduleByID gets a schedule by ID
func (s *ScheduleService) GetScheduleByID(ctx context.Context, id string) (*models.Schedule, error) {
	return s.repo.GetScheduleByID(ctx, id)
}

// GetSchedulesRequest represents get schedules request
type GetSchedulesRequest struct {
	DosenID   string `form:"dosen_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

// GetSchedules gets schedules
func (s *ScheduleService) GetSchedules(ctx context.Context, req GetSchedulesRequest) ([]models.Schedule, int64, error) {
	var startDate, endDate *time.Time
	var dosenID *string

	if req.DosenID != "" {
		dosenID = &req.DosenID
	}

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = &t
		}
	}

	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
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

	return s.repo.GetAllSchedules(ctx, dosenID, startDate, endDate, perPage, (page-1)*perPage)
}

// GetTodaySchedules gets today's schedules
func (s *ScheduleService) GetTodaySchedules(ctx context.Context, userID string, role string) ([]models.Schedule, error) {
	return s.repo.GetTodaySchedules(ctx, userID, role)
}

// GetUpcomingSchedules gets upcoming schedules
func (s *ScheduleService) GetUpcomingSchedules(ctx context.Context, userID string, role string, limit int) ([]models.Schedule, error) {
	return s.repo.GetUpcomingSchedules(ctx, userID, role, limit)
}

// GetCalendarView gets calendar view
func (s *ScheduleService) GetCalendarView(ctx context.Context, userID string, role string, year, month int) ([]models.Schedule, error) {
	return s.repo.GetCalendarView(ctx, userID, role, year, month)
}

// UpdateScheduleRequest represents update schedule request
type UpdateScheduleRequest struct {
	CourseCode *string `json:"course_code,omitempty"`
	CourseName *string `json:"course_name,omitempty"`
	Room       *string `json:"room,omitempty"`
	StartTime  *string `json:"start_time,omitempty"`
	EndTime    *string `json:"end_time,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// UpdateSchedule updates a schedule
func (s *ScheduleService) UpdateSchedule(ctx context.Context, id string, req UpdateScheduleRequest) (*models.Schedule, error) {
	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("schedule", id)
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

// DeleteSchedule deletes a schedule
func (s *ScheduleService) DeleteSchedule(ctx context.Context, id string) error {
	_, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("schedule", id)
	}
	return s.repo.DeleteSchedule(ctx, id)
}

