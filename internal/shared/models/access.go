package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessLog represents an access log entry
type AccessLog struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	GateID      string    `gorm:"type:varchar(100)" json:"gate_id"`
	AccessType  string    `gorm:"type:varchar(20);not null" json:"access_type"` // entry, exit
	IsAllowed   bool      `gorm:"default:true" json:"is_allowed"`
	Reason      string    `gorm:"type:text" json:"reason,omitempty"`
	QRCodeID    *string   `gorm:"type:uuid" json:"qr_code_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (AccessLog) TableName() string {
	return "access_logs"
}

// BeforeCreate hook
func (a *AccessLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// AccessPermission represents user access permissions
type AccessPermission struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	GateID      string    `gorm:"type:varchar(100);not null;index" json:"gate_id"`
	IsAllowed   bool      `gorm:"default:true" json:"is_allowed"`
	ValidFrom   *time.Time `json:"valid_from,omitempty"`
	ValidUntil  *time.Time `json:"valid_until,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (AccessPermission) TableName() string {
	return "access_permissions"
}

// BeforeCreate hook
func (a *AccessPermission) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

