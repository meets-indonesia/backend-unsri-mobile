package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Geofence represents a geofence area
type Geofence struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Latitude    float64   `gorm:"not null" json:"latitude"`
	Longitude   float64   `gorm:"not null" json:"longitude"`
	Radius      float64   `gorm:"not null" json:"radius"` // in meters
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (Geofence) TableName() string {
	return "geofences"
}

// BeforeCreate hook
func (g *Geofence) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

// LocationHistory represents location history for tap in/out
type LocationHistory struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string    `gorm:"type:varchar(20);not null" json:"type"` // tap_in, tap_out
	Latitude    float64   `gorm:"not null" json:"latitude"`
	Longitude   float64   `gorm:"not null" json:"longitude"`
	GeofenceID  *string   `gorm:"type:uuid;index" json:"geofence_id,omitempty"`
	IsValid     bool      `gorm:"default:true" json:"is_valid"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (LocationHistory) TableName() string {
	return "location_history"
}

// BeforeCreate hook
func (l *LocationHistory) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

