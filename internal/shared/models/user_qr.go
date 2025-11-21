package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAccessQR represents a unique QR code for user gate access
type UserAccessQR struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	QRToken   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"qr_token"` // Unique token for QR
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

