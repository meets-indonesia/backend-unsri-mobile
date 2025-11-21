package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File represents a stored file
type File struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName string   `gorm:"type:varchar(255);not null" json:"original_name"`
	FileType    string    `gorm:"type:varchar(50)" json:"file_type"` // image, document, avatar, etc.
	MimeType    string    `gorm:"type:varchar(100)" json:"mime_type"`
	Size        int64     `gorm:"not null" json:"size"` // in bytes
	Path        string    `gorm:"type:varchar(500);not null" json:"path"`
	URL         string    `gorm:"type:varchar(500)" json:"url"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (File) TableName() string {
	return "files"
}

// BeforeCreate hook
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

