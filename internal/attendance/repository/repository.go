package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// AttendanceRepository handles attendance data operations
type AttendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new attendance repository
func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// CreateAttendance creates a new attendance record
func (r *AttendanceRepository) CreateAttendance(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

// GetAttendanceByID gets an attendance by ID
func (r *AttendanceRepository) GetAttendanceByID(ctx context.Context, id string) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attendance not found")
		}
		return nil, err
	}
	return &attendance, nil
}

// GetAttendancesByUserID gets attendances by user ID
func (r *AttendanceRepository) GetAttendancesByUserID(ctx context.Context, userID string, startDate, endDate *time.Time, limit, offset int) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Attendance{}).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Order("date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// GetAttendancesByScheduleID gets attendances by schedule ID
func (r *AttendanceRepository) GetAttendancesByScheduleID(ctx context.Context, scheduleID string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := r.db.WithContext(ctx).Preload("User").Where("schedule_id = ?", scheduleID).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

// UpdateAttendance updates an attendance record
func (r *AttendanceRepository) UpdateAttendance(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

// CreateSession creates a new attendance session
func (r *AttendanceRepository) CreateSession(ctx context.Context, session *models.AttendanceSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByID gets a session by ID
func (r *AttendanceRepository) GetSessionByID(ctx context.Context, id string) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetActiveSessionByScheduleID gets active session by schedule ID
func (r *AttendanceRepository) GetActiveSessionByScheduleID(ctx context.Context, scheduleID string) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	if err := r.db.WithContext(ctx).
		Where("schedule_id = ? AND is_active = ? AND expires_at > ?", scheduleID, true, time.Now()).
		Order("created_at DESC").
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active session found")
		}
		return nil, err
	}
	return &session, nil
}

// UpdateSession updates an attendance session
func (r *AttendanceRepository) UpdateSession(ctx context.Context, session *models.AttendanceSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// CreateSchedule creates a new schedule
func (r *AttendanceRepository) CreateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetScheduleByID gets a schedule by ID
func (r *AttendanceRepository) GetScheduleByID(ctx context.Context, id string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}

// GetSchedulesByDosenID gets schedules by dosen ID
func (r *AttendanceRepository) GetSchedulesByDosenID(ctx context.Context, dosenID string, date *time.Time) ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := r.db.WithContext(ctx).Where("dosen_id = ? AND is_active = ?", dosenID, true)

	if date != nil {
		query = query.Where("date = ?", date.Format("2006-01-02"))
	}

	if err := query.Order("start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// CheckAttendanceExists checks if attendance already exists for user and date
func (r *AttendanceRepository) CheckAttendanceExists(ctx context.Context, userID string, date time.Time, scheduleID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02"))

	if scheduleID != nil {
		query = query.Where("schedule_id = ?", *scheduleID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetAttendancesByCourseID gets attendances by course ID (via schedule)
func (r *AttendanceRepository) GetAttendancesByCourseID(ctx context.Context, courseID string, startDate, endDate *time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.WithContext(ctx).Preload("User").
		Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
		Where("schedules.course_id = ?", courseID)

	if startDate != nil {
		query = query.Where("attendances.date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("attendances.date <= ?", endDate)
	}

	if err := query.Order("attendances.date DESC").Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

// GetAttendancesByStudentID gets all attendances for a student
func (r *AttendanceRepository) GetAttendancesByStudentID(ctx context.Context, studentID string, startDate, endDate *time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", studentID)

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Order("date DESC").Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

// GetAttendanceStatistics gets attendance statistics for a user
func (r *AttendanceRepository) GetAttendanceStatistics(ctx context.Context, userID string, startDate, endDate *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := r.db.WithContext(ctx).Model(&models.Attendance{}).Where("user_id = ?", userID)
	
	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	// Total attendances
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := query.Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}
	stats["by_status"] = statusMap

	// Count by type
	var typeCounts []struct {
		Type  string
		Count int64
	}
	if err := query.Select("type, COUNT(*) as count").
		Group("type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	typeMap := make(map[string]int64)
	for _, tc := range typeCounts {
		typeMap[tc.Type] = tc.Count
	}
	stats["by_type"] = typeMap

	// Attendance rate
	if total > 0 {
		hadir := statusMap["hadir"]
		stats["attendance_rate"] = float64(hadir) / float64(total) * 100
	} else {
		stats["attendance_rate"] = 0.0
	}

	return stats, nil
}

// GetTodaySchedules gets today's schedules for a user
func (r *AttendanceRepository) GetTodaySchedules(ctx context.Context, userID string, role string) ([]models.Schedule, error) {
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
func (r *AttendanceRepository) GetUpcomingSchedules(ctx context.Context, userID string, role string, limit int) ([]models.Schedule, error) {
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

// GetCurrentTapInStatus gets current tap in status for campus attendance
func (r *AttendanceRepository) GetCurrentTapInStatus(ctx context.Context, userID string) (*models.Attendance, error) {
	today := time.Now().Format("2006-01-02")
	var attendance models.Attendance
	
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ? AND type = ? AND check_in_time IS NOT NULL AND check_out_time IS NULL", 
			userID, today, models.AttendanceTypeKampus).
		Order("check_in_time DESC").
		First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active tap in
		}
		return nil, err
	}
	
	return &attendance, nil
}

// GetAllSchedules gets all schedules with filters
func (r *AttendanceRepository) GetAllSchedules(ctx context.Context, dosenID *string, startDate, endDate *time.Time, limit, offset int) ([]models.Schedule, int64, error) {
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

// UpdateSchedule updates a schedule
func (r *AttendanceRepository) UpdateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// DeleteSchedule soft deletes a schedule
func (r *AttendanceRepository) DeleteSchedule(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, "id = ?", id).Error
}

