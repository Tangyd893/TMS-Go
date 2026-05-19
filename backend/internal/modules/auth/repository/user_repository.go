package repository

import (
	"context"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login_at", loginTime).Error
}
