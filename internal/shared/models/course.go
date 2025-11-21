package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Course represents a course/mata kuliah
type Course struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	NameEn      string    `gorm:"type:varchar(255)" json:"name_en"`
	Credits     int       `gorm:"not null;default:0" json:"credits"`
	Semester    int       `json:"semester"`
	Prodi       string    `gorm:"type:varchar(255)" json:"prodi"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Classes []Class `gorm:"foreignKey:CourseID" json:"classes,omitempty"`
}

// TableName specifies the table name
func (Course) TableName() string {
	return "courses"
}

// BeforeCreate hook
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Class represents a class for a course
type Class struct {
	ID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CourseID        string    `gorm:"type:uuid;not null;index" json:"course_id"`
	ClassCode       string    `gorm:"type:varchar(50);not null" json:"class_code"`
	ClassName       string    `gorm:"type:varchar(255)" json:"class_name"`
	Semester        string    `gorm:"type:varchar(20);not null" json:"semester"`
	AcademicYear    string    `gorm:"type:varchar(20)" json:"academic_year"`
	Capacity        int       `gorm:"default:0" json:"capacity"`
	Enrolled        int       `gorm:"default:0" json:"enrolled"`
	DosenID         string    `gorm:"type:uuid;not null;index" json:"dosen_id"`
	AssistantDosenID *string  `gorm:"type:uuid;index" json:"assistant_dosen_id"`
	Room            string    `gorm:"type:varchar(100)" json:"room"`
	DayOfWeek       int       `gorm:"not null" json:"day_of_week"`
	StartTime       time.Time `gorm:"not null" json:"start_time"`
	EndTime         time.Time `gorm:"not null" json:"end_time"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course          Course       `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Enrollments    []Enrollment `gorm:"foreignKey:ClassID" json:"enrollments,omitempty"`
}

// TableName specifies the table name
func (Class) TableName() string {
	return "classes"
}

// BeforeCreate hook
func (c *Class) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Enrollment represents student enrollment in a class (KRS)
type Enrollment struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID     string    `gorm:"type:uuid;not null;index" json:"student_id"`
	ClassID       string    `gorm:"type:uuid;not null;index" json:"class_id"`
	EnrollmentDate time.Time `gorm:"not null" json:"enrollment_date"`
	Status        string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, completed, dropped, failed
	Grade         string    `gorm:"type:varchar(5)" json:"grade"` // A, B, C, D, E
	Score         float64   `gorm:"type:decimal(5,2)" json:"score"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Student User  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Class   Class `gorm:"foreignKey:ClassID" json:"class,omitempty"`
}

// TableName specifies the table name
func (Enrollment) TableName() string {
	return "enrollments"
}

// BeforeCreate hook
func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

