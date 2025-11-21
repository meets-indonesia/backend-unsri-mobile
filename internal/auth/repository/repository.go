package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// AuthRepository handles authentication data operations
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindByEmail finds a user by email
func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *AuthRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Mahasiswa").Preload("Dosen").Preload("Staff").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *AuthRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates a user
func (r *AuthRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// FindByNIM finds a mahasiswa by NIM
func (r *AuthRepository) FindByNIM(ctx context.Context, nim string) (*models.Mahasiswa, error) {
	var mahasiswa models.Mahasiswa
	if err := r.db.WithContext(ctx).Where("nim = ?", nim).First(&mahasiswa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("mahasiswa not found")
		}
		return nil, err
	}
	return &mahasiswa, nil
}

// FindByNIP finds a dosen or staff by NIP
func (r *AuthRepository) FindByNIP(ctx context.Context, nip string, role models.UserRole) (interface{}, error) {
	if role == models.RoleDosen {
		var dosen models.Dosen
		if err := r.db.WithContext(ctx).Where("nip = ?", nip).First(&dosen).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("dosen not found")
			}
			return nil, err
		}
		return &dosen, nil
	} else if role == models.RoleStaff {
		var staff models.Staff
		if err := r.db.WithContext(ctx).Where("nip = ?", nip).First(&staff).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("staff not found")
			}
			return nil, err
		}
		return &staff, nil
	}
	return nil, errors.New("invalid role")
}

// CreateMahasiswa creates a new mahasiswa
func (r *AuthRepository) CreateMahasiswa(ctx context.Context, mahasiswa *models.Mahasiswa) error {
	return r.db.WithContext(ctx).Create(mahasiswa).Error
}

// CreateDosen creates a new dosen
func (r *AuthRepository) CreateDosen(ctx context.Context, dosen *models.Dosen) error {
	return r.db.WithContext(ctx).Create(dosen).Error
}

// CreateStaff creates a new staff
func (r *AuthRepository) CreateStaff(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

