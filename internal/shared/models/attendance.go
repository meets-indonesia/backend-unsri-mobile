package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttendanceStatus represents attendance status
type AttendanceStatus string

const (
	StatusHadir    AttendanceStatus = "hadir"
	StatusIzin     AttendanceStatus = "izin"
	StatusSakit    AttendanceStatus = "sakit"
	StatusAlpa     AttendanceStatus = "alpa"
	StatusTerlambat AttendanceStatus = "terlambat"
)

// AttendanceType represents the type of attendance
type AttendanceType string

const (
	AttendanceTypeKelas AttendanceType = "kelas" // Class attendance
	AttendanceTypeKampus AttendanceType = "kampus" // Campus attendance (tap in/out)
)

// AttendanceSession represents an attendance session (for QR code generation)
type AttendanceSession struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ScheduleID *string   `gorm:"type:uuid" json:"schedule_id"` // Optional, for class attendance
	CreatedBy  string    `gorm:"type:uuid;not null" json:"created_by"` // User ID who created the session
	Type       AttendanceType `gorm:"type:varchar(20);not null" json:"type"`
	QRCode     string    `gorm:"type:text" json:"qr_code"` // QR code data
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Attendances []Attendance `gorm:"foreignKey:SessionID" json:"attendances,omitempty"`
}

// TableName specifies the table name
func (AttendanceSession) TableName() string {
	return "attendance_sessions"
}

// BeforeCreate hook to generate UUID
func (a *AttendanceSession) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// Attendance represents an attendance record
type Attendance struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	SessionID *string   `gorm:"type:uuid;index" json:"session_id"` // Optional, for QR-based attendance
	ScheduleID *string  `gorm:"type:uuid;index" json:"schedule_id"` // Optional, for class attendance
	Type      AttendanceType `gorm:"type:varchar(20);not null" json:"type"`
	Status    AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	Date      time.Time `gorm:"not null;index" json:"date"`
	CheckInTime *time.Time `json:"check_in_time"` // For tap in
	CheckOutTime *time.Time `json:"check_out_time"` // For tap out
	Latitude  *float64 `json:"latitude"` // Location latitude
	Longitude *float64 `json:"longitude"` // Location longitude
	Notes     string   `gorm:"type:text" json:"notes"`
	CreatedBy *string  `gorm:"type:uuid" json:"created_by"` // For manual entry
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (Attendance) TableName() string {
	return "attendances"
}

// BeforeCreate hook to generate UUID
func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// Schedule represents a class schedule (expandable for academic system)
type Schedule struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CourseID  *string   `gorm:"type:uuid;index" json:"course_id"` // For future expansion
	CourseCode string   `gorm:"type:varchar(50)" json:"course_code"` // Temporary, until course service is ready
	CourseName string   `gorm:"type:varchar(255)" json:"course_name"` // Temporary
	DosenID   string    `gorm:"type:uuid;not null;index" json:"dosen_id"`
	Room      string    `gorm:"type:varchar(100)" json:"room"`
	DayOfWeek int       `gorm:"not null" json:"day_of_week"` // 0=Sunday, 1=Monday, etc.
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Date      time.Time `gorm:"not null;index" json:"date"` // Specific date for this schedule
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Attendances []Attendance `gorm:"foreignKey:ScheduleID" json:"attendances,omitempty"`
}

// TableName specifies the table name
func (Schedule) TableName() string {
	return "schedules"
}

// BeforeCreate hook to generate UUID
func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

