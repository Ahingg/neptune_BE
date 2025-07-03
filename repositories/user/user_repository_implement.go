package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	model "neptune/backend/models/user"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetUserByUsername retrieves a user by their username
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when user doesn't exist
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	// Generate UUID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.WithContext(ctx).Create(user).Error
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// GetUserByID retrieves a user by their ID
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when user doesn't exist
		}
		return nil, err
	}

	return &user, nil
}
