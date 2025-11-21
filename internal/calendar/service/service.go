package service

import (
	"context"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/calendar/repository"
)

// CalendarService handles calendar business logic
type CalendarService struct {
	repo *repository.CalendarRepository
}

// NewCalendarService creates a new calendar service
func NewCalendarService(repo *repository.CalendarRepository) *CalendarService {
	return &CalendarService{repo: repo}
}

// CreateEventRequest represents create event request
type CreateEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description,omitempty"`
	EventType   string `json:"event_type,omitempty"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	IsAllDay    bool   `json:"is_all_day,omitempty"`
	Location    string `json:"location,omitempty"`
}

// CreateEvent creates a new academic event
func (s *CalendarService) CreateEvent(ctx context.Context, createdBy string, req CreateEventRequest) (*models.AcademicEvent, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_date format, use RFC3339")
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_date format, use RFC3339")
	}

	event := &models.AcademicEvent{
		Title:       req.Title,
		Description: req.Description,
		EventType:   req.EventType,
		StartDate:   startDate,
		EndDate:     endDate,
		IsAllDay:    req.IsAllDay,
		Location:    req.Location,
		CreatedBy:   createdBy,
		IsActive:    true,
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, apperrors.NewInternalError("failed to create event", err)
	}

	return event, nil
}

// GetEventByID gets an event by ID
func (s *CalendarService) GetEventByID(ctx context.Context, id string) (*models.AcademicEvent, error) {
	return s.repo.GetEventByID(ctx, id)
}

// GetEventsRequest represents get events request
type GetEventsRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	EventType string `form:"event_type"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

// GetEvents gets all events
func (s *CalendarService) GetEvents(ctx context.Context, req GetEventsRequest) ([]models.AcademicEvent, int64, error) {
	var startDate, endDate *time.Time
	var eventType *string

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

	if req.EventType != "" {
		eventType = &req.EventType
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.GetAllEvents(ctx, startDate, endDate, eventType, perPage, (page-1)*perPage)
}

// GetEventsByMonth gets events for a month
func (s *CalendarService) GetEventsByMonth(ctx context.Context, year, month int) ([]models.AcademicEvent, error) {
	return s.repo.GetEventsByMonth(ctx, year, month)
}

// GetUpcomingEvents gets upcoming events
func (s *CalendarService) GetUpcomingEvents(ctx context.Context, limit int) ([]models.AcademicEvent, error) {
	return s.repo.GetUpcomingEvents(ctx, limit)
}

// UpdateEventRequest represents update event request
type UpdateEventRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	EventType   *string `json:"event_type,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
	IsAllDay    *bool   `json:"is_all_day,omitempty"`
	Location    *string `json:"location,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdateEvent updates an event
func (s *CalendarService) UpdateEvent(ctx context.Context, id string, req UpdateEventRequest) (*models.AcademicEvent, error) {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("event", id)
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.EventType != nil {
		event.EventType = *req.EventType
	}
	if req.StartDate != nil {
		startDate, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid start_date format")
		}
		event.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid end_date format")
		}
		event.EndDate = endDate
	}
	if req.IsAllDay != nil {
		event.IsAllDay = *req.IsAllDay
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.IsActive != nil {
		event.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateEvent(ctx, event); err != nil {
		return nil, apperrors.NewInternalError("failed to update event", err)
	}

	return event, nil
}

// DeleteEvent deletes an event
func (s *CalendarService) DeleteEvent(ctx context.Context, id string) error {
	_, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("event", id)
	}
	return s.repo.DeleteEvent(ctx, id)
}

