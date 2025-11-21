package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transcript represents student transcript
type Transcript struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID string    `gorm:"type:uuid;not null;index" json:"student_id"`
	Semester  string    `gorm:"type:varchar(20);not null" json:"semester"`
	IPK       float64   `gorm:"type:decimal(3,2)" json:"ipk"`
	Data      string    `gorm:"type:jsonb" json:"data"` // Course grades as JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name
func (Transcript) TableName() string {
	return "transcripts"
}

// BeforeCreate hook
func (t *Transcript) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// KRS represents Kartu Rencana Studi
type KRS struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID string    `gorm:"type:uuid;not null;index" json:"student_id"`
	Semester  string    `gorm:"type:varchar(20);not null" json:"semester"`
	Status    string    `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, approved, rejected
	Data      string    `gorm:"type:jsonb" json:"data"` // Courses as JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name
func (KRS) TableName() string {
	return "krss"
}

// BeforeCreate hook
func (k *KRS) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = uuid.New().String()
	}
	return nil
}

// Bimbingan represents guidance/consultation record
type Bimbingan struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID string    `gorm:"type:uuid;not null;index" json:"student_id"`
	DosenID   string    `gorm:"type:uuid;not null;index" json:"dosen_id"`
	Topic     string    `gorm:"type:varchar(255)" json:"topic"`
	Notes     string    `gorm:"type:text" json:"notes"`
	Date      time.Time `gorm:"not null" json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Dosen   User `gorm:"foreignKey:DosenID" json:"dosen,omitempty"`
}

// TableName specifies the table name
func (Bimbingan) TableName() string {
	return "bimbingans"
}

// BeforeCreate hook
func (b *Bimbingan) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

