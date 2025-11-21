package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AcademicEvent represents an academic calendar event
type AcademicEvent struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	EventType   string    `gorm:"type:varchar(50)" json:"event_type"` // exam, holiday, registration, etc.
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	IsAllDay    bool      `gorm:"default:false" json:"is_all_day"`
	Location    string    `gorm:"type:varchar(255)" json:"location"`
	CreatedBy   string    `gorm:"type:uuid;not null" json:"created_by"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (AcademicEvent) TableName() string {
	return "academic_events"
}

// BeforeCreate hook
func (e *AcademicEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

