package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// ReportRepository handles report data operations
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetAttendanceReport gets attendance report for a period
func (r *ReportRepository) GetAttendanceReport(ctx context.Context, studentID *string, courseID *string, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	if studentID != nil {
		query = query.Where("user_id = ?", *studentID)
	}

	if courseID != nil {
		query = query.Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
			Where("schedules.course_id = ?", *courseID)
	}

	if err := query.Order("created_at DESC").Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

// GetAttendanceSummary gets attendance summary
func (r *ReportRepository) GetAttendanceSummary(ctx context.Context, studentID *string, courseID *string, startDate, endDate time.Time) (map[string]interface{}, error) {
	var total, present, absent, late, excused int64

	query := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	if studentID != nil {
		query = query.Where("user_id = ?", *studentID)
	}

	if courseID != nil {
		query = query.Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
			Where("schedules.course_id = ?", *courseID)
	}

	query.Count(&total)
	query.Where("status = ?", "present").Count(&present)
	query.Where("status = ?", "absent").Count(&absent)
	query.Where("status = ?", "late").Count(&late)
	query.Where("status = ?", "excused").Count(&excused)

	attendanceRate := float64(0)
	if total > 0 {
		attendanceRate = float64(present) / float64(total) * 100
	}

	return map[string]interface{}{
		"total":          total,
		"present":        present,
		"absent":         absent,
		"late":           late,
		"excused":        excused,
		"attendance_rate": attendanceRate,
	}, nil
}

// GetStudentAcademicReport gets academic report for a student
func (r *ReportRepository) GetStudentAcademicReport(ctx context.Context, studentID string, semester *string) (map[string]interface{}, error) {
	var transcript models.Transcript
	query := r.db.WithContext(ctx).Model(&models.Transcript{}).
		Where("student_id = ?", studentID)

	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}

	if err := query.Order("semester DESC").First(&transcript).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"student_id": studentID,
				"semester":   semester,
				"ipk":        0.0,
				"courses":    []interface{}{},
			}, nil
		}
		return nil, err
	}

	// Get KRS for the semester
	var krs models.KRS
	sem := semester
	if sem == nil {
		sem = &transcript.Semester
	}
	r.db.WithContext(ctx).Where("student_id = ? AND semester = ?", studentID, *sem).First(&krs)

	return map[string]interface{}{
		"student_id": studentID,
		"semester":   transcript.Semester,
		"ipk":        transcript.IPK,
		"transcript": transcript,
		"krs":        krs,
	}, nil
}

// GetCourseReport gets report for a course
func (r *ReportRepository) GetCourseReport(ctx context.Context, courseID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Where("id = ?", courseID).First(&course).Error; err != nil {
		return nil, err
	}

	// Get schedules for this course
	var schedules []models.Schedule
	r.db.WithContext(ctx).Where("course_id = ?", courseID).Find(&schedules)

	// Get attendance stats
	var totalAttendance, totalStudents int64
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
		Where("schedules.course_id = ? AND attendances.created_at >= ? AND attendances.created_at <= ?",
			courseID, startDate, endDate).
		Count(&totalAttendance)

	// Count distinct students who have attendance for this course
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
		Where("schedules.course_id = ?", courseID).
		Distinct("attendances.user_id").
		Count(&totalStudents)

	return map[string]interface{}{
		"course":          course,
		"schedules":       schedules,
		"total_attendance": totalAttendance,
		"total_students":   totalStudents,
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}, nil
}

// GetDailyReport gets daily report
func (r *ReportRepository) GetDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalAttendance, totalTapIn, totalTapOut int64

	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&totalAttendance)

	r.db.WithContext(ctx).Model(&models.LocationHistory{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "tap_in", startOfDay, endOfDay).
		Count(&totalTapIn)

	r.db.WithContext(ctx).Model(&models.LocationHistory{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "tap_out", startOfDay, endOfDay).
		Count(&totalTapOut)

	return map[string]interface{}{
		"date":             date.Format("2006-01-02"),
		"total_attendance": totalAttendance,
		"total_tap_in":     totalTapIn,
		"total_tap_out":    totalTapOut,
	}, nil
}

