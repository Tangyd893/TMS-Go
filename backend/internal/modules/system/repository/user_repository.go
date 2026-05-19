package repository

import (
	"context"
	"errors"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	List(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User, roleIDs []string) error
	Update(ctx context.Context, user *model.User, roleIDs []string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Preload("Roles")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user.ID.String(), roleIDs)
	})
}

func (r *userRepository) Update(ctx context.Context, user *model.User, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, err := r.FindByID(ctx, user.ID)
		if err != nil {
			return err
		}

		updates := map[string]any{
			"real_name":  user.RealName,
			"phone":      user.Phone,
			"email":      user.Email,
			"status":     user.Status,
			"updated_at": gorm.Expr("now()"),
		}
		if user.Username != "" {
			updates["username"] = user.Username
		}
		_ = existing

		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user.ID.String(), roleIDs)
	})
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) syncRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	for _, roleID := range roleIDs {
		ur := &model.UserRole{
			UserID: uuid.MustParse(userID),
			RoleID: uuid.MustParse(roleID),
		}
		if err := tx.Create(ur).Error; err != nil {
			return err
		}
	}
	return nil
}

var ErrRecordNotFound = errors.New("record not found")
