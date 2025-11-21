package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// CalendarRepository handles calendar data operations
type CalendarRepository struct {
	db *gorm.DB
}

// NewCalendarRepository creates a new calendar repository
func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

// CreateEvent creates a new academic event
func (r *CalendarRepository) CreateEvent(ctx context.Context, event *models.AcademicEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetEventByID gets an event by ID
func (r *CalendarRepository) GetEventByID(ctx context.Context, id string) (*models.AcademicEvent, error) {
	var event models.AcademicEvent
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

// GetAllEvents gets all events with filters
func (r *CalendarRepository) GetAllEvents(ctx context.Context, startDate, endDate *time.Time, eventType *string, limit, offset int) ([]models.AcademicEvent, int64, error) {
	var events []models.AcademicEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AcademicEvent{}).Where("is_active = ?", true)

	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", endDate)
	}
	if eventType != nil {
		query = query.Where("event_type = ?", *eventType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("start_date ASC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetEventsByMonth gets events for a specific month
func (r *CalendarRepository) GetEventsByMonth(ctx context.Context, year, month int) ([]models.AcademicEvent, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var events []models.AcademicEvent
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND start_date >= ? AND start_date <= ?", true, startDate, endDate).
		Order("start_date ASC").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// GetUpcomingEvents gets upcoming events
func (r *CalendarRepository) GetUpcomingEvents(ctx context.Context, limit int) ([]models.AcademicEvent, error) {
	now := time.Now()
	var events []models.AcademicEvent

	query := r.db.WithContext(ctx).
		Where("is_active = ? AND start_date >= ?", true, now).
		Order("start_date ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates an event
func (r *CalendarRepository) UpdateEvent(ctx context.Context, event *models.AcademicEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

// DeleteEvent soft deletes an event
func (r *CalendarRepository) DeleteEvent(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AcademicEvent{}, "id = ?", id).Error
}

