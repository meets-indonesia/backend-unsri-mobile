package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType represents notification type
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeSuccess NotificationType = "success"
)

// Notification represents a notification
type Notification struct {
	ID        string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string           `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string           `gorm:"type:varchar(255);not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	Data      string           `gorm:"type:jsonb" json:"data,omitempty"` // Additional data as JSON
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// DeviceToken represents FCM device token
type DeviceToken struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	Platform  string    `gorm:"type:varchar(20)" json:"platform"` // ios, android, web
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (DeviceToken) TableName() string {
	return "device_tokens"
}

// BeforeCreate hook
func (d *DeviceToken) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

