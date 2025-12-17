package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAccessQR represents a unique QR code for user gate access
type UserAccessQR struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	SessionID string         `gorm:"type:uuid;not null;uniqueIndex" json:"session_id"`       // Unique session ID per QR generation
	QRToken   string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"qr_token"` // Unique token for QR (legacy, kept for backward compatibility)
	IsActive  bool           `gorm:"default:true" json:"is_active"`                          // true = tap-in (masuk), false = tap-out (keluar)
	ExpiresAt *time.Time     `gorm:"type:timestamp" json:"expires_at,omitempty"`             // Set to now when tap-out, null when active
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (UserAccessQR) TableName() string {
	return "user_access_qrs"
}

// BeforeCreate hook
func (u *UserAccessQR) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
