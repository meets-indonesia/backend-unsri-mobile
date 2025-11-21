package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BroadcastType represents broadcast type
type BroadcastType string

const (
	BroadcastTypeGeneral BroadcastType = "general"
	BroadcastTypeClass   BroadcastType = "class"
	BroadcastTypeCampus  BroadcastType = "campus"
)

// Broadcast represents a broadcast/pengumuman
type Broadcast struct {
	ID          string        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string        `gorm:"type:varchar(255);not null" json:"title"`
	Content     string        `gorm:"type:text;not null" json:"content"`
	Type        BroadcastType `gorm:"type:varchar(20);not null" json:"type"`
	Priority    string        `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent
	CreatedBy   string        `gorm:"type:uuid;not null" json:"created_by"`
	ClassID     *string       `gorm:"type:uuid;index" json:"class_id,omitempty"` // For class-specific broadcasts
	IsPublished bool          `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time    `json:"published_at,omitempty"`
	ScheduledAt *time.Time    `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time    `json:"expires_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Audiences []BroadcastAudience `gorm:"foreignKey:BroadcastID" json:"audiences,omitempty"`
}

// TableName specifies the table name
func (Broadcast) TableName() string {
	return "broadcasts"
}

// BeforeCreate hook
func (b *Broadcast) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BroadcastAudience represents broadcast target audience
type BroadcastAudience struct {
	ID         string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BroadcastID string `gorm:"type:uuid;not null;index" json:"broadcast_id"`
	UserID     *string `gorm:"type:uuid;index" json:"user_id,omitempty"` // Specific user
	Role       *string `gorm:"type:varchar(20)" json:"role,omitempty"` // All users with this role
	Prodi      *string `gorm:"type:varchar(255)" json:"prodi,omitempty"` // All users in this prodi
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name
func (BroadcastAudience) TableName() string {
	return "broadcast_audiences"
}

// BeforeCreate hook
func (b *BroadcastAudience) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

