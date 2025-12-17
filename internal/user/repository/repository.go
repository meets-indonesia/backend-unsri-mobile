package repository

import (
	"context"
	"errors"

	"unsri-backend/internal/shared/models"

	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID gets a user by ID with role-specific data
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Mahasiswa").
		Preload("Dosen").
		Preload("Staff").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetMahasiswaByNIM gets mahasiswa by NIM
func (r *UserRepository) GetMahasiswaByNIM(ctx context.Context, nim string) (*models.Mahasiswa, error) {
	var mahasiswa models.Mahasiswa
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("nim = ?", nim).
		First(&mahasiswa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("mahasiswa not found")
		}
		return nil, err
	}
	return &mahasiswa, nil
}

// GetDosenByNIP gets dosen by NIP
func (r *UserRepository) GetDosenByNIP(ctx context.Context, nip string) (*models.Dosen, error) {
	var dosen models.Dosen
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("nip = ?", nip).
		First(&dosen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dosen not found")
		}
		return nil, err
	}
	return &dosen, nil
}

// GetStaffByNIP gets staff by NIP
func (r *UserRepository) GetStaffByNIP(ctx context.Context, nip string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("nip = ?", nip).
		First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, err
	}
	return &staff, nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateMahasiswa updates mahasiswa data
func (r *UserRepository) UpdateMahasiswa(ctx context.Context, mahasiswa *models.Mahasiswa) error {
	return r.db.WithContext(ctx).Save(mahasiswa).Error
}

// UpdateDosen updates dosen data
func (r *UserRepository) UpdateDosen(ctx context.Context, dosen *models.Dosen) error {
	return r.db.WithContext(ctx).Save(dosen).Error
}

// UpdateStaff updates staff data
func (r *UserRepository) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Save(staff).Error
}

// SearchUsers searches users by query
func (r *UserRepository) SearchUsers(ctx context.Context, query string, role *string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&models.User{})

	if query != "" {
		searchPattern := "%" + query + "%"
		dbQuery = dbQuery.Where("email ILIKE ?", searchPattern)
	}

	if role != nil {
		dbQuery = dbQuery.Where("role = ?", *role)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.
		Preload("Mahasiswa").
		Preload("Dosen").
		Preload("Staff").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListUsers lists all users with pagination and filters
func (r *UserRepository) ListUsers(ctx context.Context, role *string, isActive *bool, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&models.User{})

	if role != nil {
		dbQuery = dbQuery.Where("role = ?", *role)
	}

	if isActive != nil {
		dbQuery = dbQuery.Where("is_active = ?", *isActive)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.
		Preload("Mahasiswa").
		Preload("Dosen").
		Preload("Staff").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ActivateUser activates a user
func (r *UserRepository) ActivateUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeactivateUser deactivates a user
func (r *UserRepository) DeactivateUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetDB returns the underlying GORM DB instance for transaction support
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

// CreateUser creates a new user (for admin)
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// CreateMahasiswa creates a new mahasiswa
func (r *UserRepository) CreateMahasiswa(ctx context.Context, mahasiswa *models.Mahasiswa) error {
	return r.db.WithContext(ctx).Create(mahasiswa).Error
}

// CreateDosen creates a new dosen
func (r *UserRepository) CreateDosen(ctx context.Context, dosen *models.Dosen) error {
	return r.db.WithContext(ctx).Create(dosen).Error
}

// CreateStaff creates a new staff
func (r *UserRepository) CreateStaff(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
