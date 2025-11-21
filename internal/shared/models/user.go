package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents user roles
type UserRole string

const (
	RoleMahasiswa UserRole = "mahasiswa"
	RoleDosen     UserRole = "dosen"
	RoleStaff     UserRole = "staff"
)

// User represents a user in the system
type User struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(20);not null" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Mahasiswa *Mahasiswa `gorm:"foreignKey:UserID" json:"mahasiswa,omitempty"`
	Dosen     *Dosen     `gorm:"foreignKey:UserID" json:"dosen,omitempty"`
	Staff     *Staff     `gorm:"foreignKey:UserID" json:"staff,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Mahasiswa represents a student
type Mahasiswa struct {
	ID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	NIM    string `gorm:"uniqueIndex;not null" json:"nim"`
	Nama   string `gorm:"not null" json:"nama"`
	Prodi  string `json:"prodi"` // Program Studi
	Angkatan int  `json:"angkatan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (Mahasiswa) TableName() string {
	return "mahasiswa"
}

// BeforeCreate hook to generate UUID
func (m *Mahasiswa) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// Dosen represents a lecturer
type Dosen struct {
	ID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	NIP    string `gorm:"uniqueIndex;not null" json:"nip"`
	Nama   string `gorm:"not null" json:"nama"`
	Prodi  string `json:"prodi"` // Program Studi
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (Dosen) TableName() string {
	return "dosen"
}

// BeforeCreate hook to generate UUID
func (d *Dosen) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// Staff represents a staff member
type Staff struct {
	ID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	NIP    string `gorm:"uniqueIndex;not null" json:"nip"`
	Nama   string `gorm:"not null" json:"nama"`
	Jabatan string `json:"jabatan"`
	Unit    string `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (Staff) TableName() string {
	return "staff"
}

// BeforeCreate hook to generate UUID
func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

