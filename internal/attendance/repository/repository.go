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

// ========== Work Attendance (HRIS) Repository Methods ==========

// CreateShiftPattern creates a new shift pattern
func (r *AttendanceRepository) CreateShiftPattern(ctx context.Context, shift *models.ShiftPattern) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

// GetShiftPatternByID gets a shift pattern by ID
func (r *AttendanceRepository) GetShiftPatternByID(ctx context.Context, id string) (*models.ShiftPattern, error) {
	var shift models.ShiftPattern
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&shift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shift pattern not found")
		}
		return nil, err
	}
	return &shift, nil
}

// GetShiftPatternByCode gets a shift pattern by code
func (r *AttendanceRepository) GetShiftPatternByCode(ctx context.Context, code string) (*models.ShiftPattern, error) {
	var shift models.ShiftPattern
	if err := r.db.WithContext(ctx).Where("shift_code = ?", code).First(&shift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shift pattern not found")
		}
		return nil, err
	}
	return &shift, nil
}

// GetAllShiftPatterns gets all shift patterns with filters
func (r *AttendanceRepository) GetAllShiftPatterns(ctx context.Context, isActive *bool, limit, offset int) ([]models.ShiftPattern, int64, error) {
	var shifts []models.ShiftPattern
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ShiftPattern{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("shift_code ASC").Find(&shifts).Error; err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

// UpdateShiftPattern updates a shift pattern
func (r *AttendanceRepository) UpdateShiftPattern(ctx context.Context, shift *models.ShiftPattern) error {
	return r.db.WithContext(ctx).Save(shift).Error
}

// DeleteShiftPattern soft deletes a shift pattern
func (r *AttendanceRepository) DeleteShiftPattern(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.ShiftPattern{}, "id = ?", id).Error
}

// CreateUserShift creates a new user shift assignment
func (r *AttendanceRepository) CreateUserShift(ctx context.Context, userShift *models.UserShift) error {
	return r.db.WithContext(ctx).Create(userShift).Error
}

// GetUserShiftByID gets a user shift by ID
func (r *AttendanceRepository) GetUserShiftByID(ctx context.Context, id string) (*models.UserShift, error) {
	var userShift models.UserShift
	if err := r.db.WithContext(ctx).Preload("User").Preload("Shift").Where("id = ?", id).First(&userShift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user shift not found")
		}
		return nil, err
	}
	return &userShift, nil
}

// GetUserShiftsByUserID gets user shifts by user ID
func (r *AttendanceRepository) GetUserShiftsByUserID(ctx context.Context, userID string, date *time.Time) ([]models.UserShift, error) {
	var userShifts []models.UserShift
	query := r.db.WithContext(ctx).Preload("Shift").Where("user_id = ? AND is_active = ?", userID, true)

	if date != nil {
		query = query.Where("effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", date, date)
	}

	if err := query.Order("effective_from DESC").Find(&userShifts).Error; err != nil {
		return nil, err
	}
	return userShifts, nil
}

// UpdateUserShift updates a user shift
func (r *AttendanceRepository) UpdateUserShift(ctx context.Context, userShift *models.UserShift) error {
	return r.db.WithContext(ctx).Save(userShift).Error
}

// DeleteUserShift soft deletes a user shift
func (r *AttendanceRepository) DeleteUserShift(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.UserShift{}, "id = ?", id).Error
}

// CreateWorkSchedule creates a new work schedule
func (r *AttendanceRepository) CreateWorkSchedule(ctx context.Context, schedule *models.WorkSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetWorkScheduleByID gets a work schedule by ID
func (r *AttendanceRepository) GetWorkScheduleByID(ctx context.Context, id string) (*models.WorkSchedule, error) {
	var schedule models.WorkSchedule
	if err := r.db.WithContext(ctx).Preload("User").Preload("Shift").Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("work schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}

// GetWorkSchedulesByUserID gets work schedules by user ID
func (r *AttendanceRepository) GetWorkSchedulesByUserID(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.WorkSchedule, error) {
	var schedules []models.WorkSchedule
	query := r.db.WithContext(ctx).Preload("Shift").Where("user_id = ? AND is_active = ?", userID, true)

	if startDate != nil {
		query = query.Where("schedule_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("schedule_date <= ?", endDate)
	}

	if err := query.Order("schedule_date ASC, start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetAllWorkSchedules gets all work schedules with filters
func (r *AttendanceRepository) GetAllWorkSchedules(ctx context.Context, userID *string, startDate, endDate *time.Time, limit, offset int) ([]models.WorkSchedule, int64, error) {
	var schedules []models.WorkSchedule
	var total int64

	query := r.db.WithContext(ctx).Model(&models.WorkSchedule{}).Where("is_active = ?", true)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if startDate != nil {
		query = query.Where("schedule_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("schedule_date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Shift").Limit(limit).Offset(offset).Order("schedule_date ASC, start_time ASC").Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

// UpdateWorkSchedule updates a work schedule
func (r *AttendanceRepository) UpdateWorkSchedule(ctx context.Context, schedule *models.WorkSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// DeleteWorkSchedule soft deletes a work schedule
func (r *AttendanceRepository) DeleteWorkSchedule(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.WorkSchedule{}, "id = ?", id).Error
}

// CreateWorkAttendanceSession creates a new work attendance session
func (r *AttendanceRepository) CreateWorkAttendanceSession(ctx context.Context, session *models.WorkAttendanceSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetWorkAttendanceSessionByID gets a work attendance session by ID
func (r *AttendanceRepository) GetWorkAttendanceSessionByID(ctx context.Context, id string) (*models.WorkAttendanceSession, error) {
	var session models.WorkAttendanceSession
	if err := r.db.WithContext(ctx).Preload("Schedule").Where("id = ?", id).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("work attendance session not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetWorkAttendanceSessionByQRCode gets a work attendance session by QR code
func (r *AttendanceRepository) GetWorkAttendanceSessionByQRCode(ctx context.Context, qrCode string) (*models.WorkAttendanceSession, error) {
	var session models.WorkAttendanceSession
	if err := r.db.WithContext(ctx).Preload("Schedule").
		Where("qr_code_data = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", qrCode, true, time.Now()).
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("work attendance session not found or expired")
		}
		return nil, err
	}
	return &session, nil
}

// UpdateWorkAttendanceSession updates a work attendance session
func (r *AttendanceRepository) UpdateWorkAttendanceSession(ctx context.Context, session *models.WorkAttendanceSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// CreateWorkAttendanceRecord creates a new work attendance record
func (r *AttendanceRepository) CreateWorkAttendanceRecord(ctx context.Context, record *models.WorkAttendanceRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetWorkAttendanceRecordByID gets a work attendance record by ID
func (r *AttendanceRepository) GetWorkAttendanceRecordByID(ctx context.Context, id string) (*models.WorkAttendanceRecord, error) {
	var record models.WorkAttendanceRecord
	if err := r.db.WithContext(ctx).Preload("User").Preload("Schedule").Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("work attendance record not found")
		}
		return nil, err
	}
	return &record, nil
}

// GetWorkAttendanceRecordsByUserID gets work attendance records by user ID
func (r *AttendanceRepository) GetWorkAttendanceRecordsByUserID(ctx context.Context, userID string, startDate, endDate *time.Time, limit, offset int) ([]models.WorkAttendanceRecord, int64, error) {
	var records []models.WorkAttendanceRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&models.WorkAttendanceRecord{}).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("DATE(recorded_at) >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("DATE(recorded_at) <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Schedule").Order("recorded_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetWorkAttendanceRecordsByScheduleID gets work attendance records by schedule ID
func (r *AttendanceRepository) GetWorkAttendanceRecordsByScheduleID(ctx context.Context, scheduleID string) ([]models.WorkAttendanceRecord, error) {
	var records []models.WorkAttendanceRecord
	if err := r.db.WithContext(ctx).Preload("User").
		Where("schedule_id = ?", scheduleID).
		Order("recorded_at ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetTodayWorkAttendanceRecord gets today's work attendance record for check-in/out
func (r *AttendanceRepository) GetTodayWorkAttendanceRecord(ctx context.Context, userID string, attendanceType string) (*models.WorkAttendanceRecord, error) {
	today := time.Now().Format("2006-01-02")
	var record models.WorkAttendanceRecord

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND DATE(recorded_at) = ? AND attendance_type = ?", userID, today, attendanceType).
		Order("recorded_at DESC").
		First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, err
	}

	return &record, nil
}

// UpdateWorkAttendanceRecord updates a work attendance record
func (r *AttendanceRepository) UpdateWorkAttendanceRecord(ctx context.Context, record *models.WorkAttendanceRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

// GetFileByID gets a file by ID
func (r *AttendanceRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	var file models.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}
