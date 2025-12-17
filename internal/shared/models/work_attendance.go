package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkAttendanceStatus represents work attendance status
type WorkAttendanceStatus string

const (
	StatusCheckIn   WorkAttendanceStatus = "CHECK_IN"
	StatusCheckOut  WorkAttendanceStatus = "CHECK_OUT"
	StatusLateIn    WorkAttendanceStatus = "LATE_IN"
	StatusEarlyOut  WorkAttendanceStatus = "EARLY_OUT"
	StatusAbsent    WorkAttendanceStatus = "ABSENT"
	StatusOnLeave   WorkAttendanceStatus = "ON_LEAVE"
	StatusSickLeave WorkAttendanceStatus = "SICK_LEAVE"
)

// ShiftPattern represents a shift pattern for work schedules
type ShiftPattern struct {
	ID                   string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ShiftName            string         `gorm:"type:varchar(100);not null" json:"shift_name"`
	ShiftCode            string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"shift_code"`
	StartTime            time.Time      `gorm:"type:time;not null" json:"start_time"`
	EndTime              time.Time      `gorm:"type:time;not null" json:"end_time"`
	BreakDurationMinutes *int           `gorm:"type:integer" json:"break_duration_minutes,omitempty"`
	IsNightShift         bool           `gorm:"default:false" json:"is_night_shift"`
	IsActive             bool           `gorm:"default:true" json:"is_active"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (ShiftPattern) TableName() string {
	return "shift_patterns"
}

// BeforeCreate hook
func (s *ShiftPattern) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// UserShift represents user shift assignment
type UserShift struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         string         `gorm:"type:uuid;not null;index" json:"user_id"`
	ShiftID        string         `gorm:"type:uuid;not null;index" json:"shift_id"`
	EffectiveFrom  time.Time      `gorm:"type:date;not null" json:"effective_from"`
	EffectiveUntil *time.Time     `gorm:"type:date" json:"effective_until,omitempty"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (excluded from JSON response)
	User  User         `gorm:"foreignKey:UserID" json:"-"`
	Shift ShiftPattern `gorm:"foreignKey:ShiftID" json:"-"`
}

// TableName specifies the table name
func (UserShift) TableName() string {
	return "user_shifts"
}

// BeforeCreate hook
func (u *UserShift) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// WorkSchedule represents a work schedule for a user
type WorkSchedule struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string         `gorm:"type:uuid;not null;index" json:"user_id"`
	ScheduleDate time.Time      `gorm:"type:date;not null" json:"schedule_date"`
	DayOfWeek    *int           `gorm:"type:integer" json:"day_of_week,omitempty"` // 0=Sunday, 1=Monday, etc.
	ShiftID      *string        `gorm:"type:uuid;index" json:"shift_id,omitempty"`
	StartTime    time.Time      `gorm:"type:time;not null" json:"start_time"`
	EndTime      time.Time      `gorm:"type:time;not null" json:"end_time"`
	WorkType     string         `gorm:"type:varchar(50)" json:"work_type"`
	Location     string         `gorm:"type:varchar(255)" json:"location"`
	IsHoliday    bool           `gorm:"default:false" json:"is_holiday"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (excluded from JSON response)
	User  User          `gorm:"foreignKey:UserID" json:"-"`
	Shift *ShiftPattern `gorm:"foreignKey:ShiftID" json:"-"`
}

// TableName specifies the table name
func (WorkSchedule) TableName() string {
	return "work_schedules"
}

// BeforeCreate hook
func (w *WorkSchedule) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

// WorkAttendanceSession represents a work attendance session
type WorkAttendanceSession struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ScheduleID  *string        `gorm:"type:uuid;index" json:"schedule_id,omitempty"`
	SessionDate time.Time      `gorm:"type:date;not null" json:"session_date"`
	QRCodeData  *string        `gorm:"type:varchar(255);uniqueIndex" json:"qr_code_data,omitempty"`
	ExpiresAt   *time.Time     `gorm:"type:timestamp" json:"expires_at,omitempty"`
	Status      string         `gorm:"type:varchar(20)" json:"status"` // OPEN, CLOSED, EXPIRED
	CreatedBy   *string        `gorm:"type:uuid;index" json:"created_by,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Schedule *WorkSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// TableName specifies the table name
func (WorkAttendanceSession) TableName() string {
	return "work_attendance_sessions"
}

// BeforeCreate hook
func (w *WorkAttendanceSession) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

// WorkAttendanceRecord represents a work attendance record
type WorkAttendanceRecord struct {
	ID             string               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionID      *string              `gorm:"type:uuid;index" json:"session_id,omitempty"`
	ScheduleID     *string              `gorm:"type:uuid;index" json:"schedule_id,omitempty"`
	UserID         string               `gorm:"type:uuid;not null;index" json:"user_id"`
	AttendanceType string               `gorm:"type:varchar(20);not null" json:"attendance_type"` // CHECK_IN, CHECK_OUT
	RecordedAt     time.Time            `gorm:"type:timestamp;not null" json:"recorded_at"`
	Status         WorkAttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	IsViaUNSRIWiFi *bool                `gorm:"type:boolean" json:"is_via_unsri_wifi,omitempty"`
	Latitude       *float64             `json:"latitude,omitempty"`
	Longitude      *float64             `json:"longitude,omitempty"`
	GeofenceID     *string              `gorm:"type:uuid;index" json:"geofence_id,omitempty"`
	Notes          string               `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`

	// Relations
	User     User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Schedule *WorkSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// TableName specifies the table name
func (WorkAttendanceRecord) TableName() string {
	return "work_attendance_records"
}

// BeforeCreate hook
func (w *WorkAttendanceRecord) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
