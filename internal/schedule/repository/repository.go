package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// ScheduleRepository handles schedule data operations
type ScheduleRepository struct {
	db *gorm.DB
}

// NewScheduleRepository creates a new schedule repository
func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// CreateSchedule creates a new schedule
func (r *ScheduleRepository) CreateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetScheduleByID gets a schedule by ID
func (r *ScheduleRepository) GetScheduleByID(ctx context.Context, id string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}

// GetAllSchedules gets all schedules with filters
func (r *ScheduleRepository) GetAllSchedules(ctx context.Context, dosenID *string, startDate, endDate *time.Time, limit, offset int) ([]models.Schedule, int64, error) {
	var schedules []models.Schedule
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Schedule{}).Where("is_active = ?", true)

	if dosenID != nil {
		query = query.Where("dosen_id = ?", *dosenID)
	}

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("date ASC, start_time ASC").Limit(limit).Offset(offset).Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

// GetTodaySchedules gets today's schedules
func (r *ScheduleRepository) GetTodaySchedules(ctx context.Context, userID string, role string) ([]models.Schedule, error) {
	today := time.Now()
	var schedules []models.Schedule

	query := r.db.WithContext(ctx).Where("date = ? AND is_active = ?", today.Format("2006-01-02"), true)

	if role == "dosen" {
		query = query.Where("dosen_id = ?", userID)
	}

	if err := query.Order("start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetUpcomingSchedules gets upcoming schedules
func (r *ScheduleRepository) GetUpcomingSchedules(ctx context.Context, userID string, role string, limit int) ([]models.Schedule, error) {
	today := time.Now()
	var schedules []models.Schedule

	query := r.db.WithContext(ctx).Where("date >= ? AND is_active = ?", today.Format("2006-01-02"), true)

	if role == "dosen" {
		query = query.Where("dosen_id = ?", userID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("date ASC, start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetCalendarView gets schedules for calendar view
func (r *ScheduleRepository) GetCalendarView(ctx context.Context, userID string, role string, year, month int) ([]models.Schedule, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var schedules []models.Schedule
	query := r.db.WithContext(ctx).Where("date >= ? AND date <= ? AND is_active = ?", startDate, endDate, true)

	if role == "dosen" {
		query = query.Where("dosen_id = ?", userID)
	}

	if err := query.Order("date ASC, start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// UpdateSchedule updates a schedule
func (r *ScheduleRepository) UpdateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// DeleteSchedule soft deletes a schedule
func (r *ScheduleRepository) DeleteSchedule(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, "id = ?", id).Error
}

