package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CrudRepository[T any] interface {
	List(ctx context.Context, page, pageSize int, keyword string, searchFields []string) ([]T, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, id uuid.UUID, updates map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type crudRepository[T any] struct {
	db *gorm.DB
}

func NewCrudRepository[T any](db *gorm.DB) CrudRepository[T] {
	return &crudRepository[T]{db: db}
}

func (r *crudRepository[T]) List(ctx context.Context, page, pageSize int, keyword string, searchFields []string) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))

	if keyword != "" && len(searchFields) > 0 {
		like := "%" + keyword + "%"
		var conditions []string
		var args []any
		for _, f := range searchFields {
			conditions = append(conditions, f+" LIKE ?")
			args = append(args, like)
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *crudRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *crudRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *crudRepository[T]) Update(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	updates["updated_at"] = gorm.Expr("now()")
	return r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func (r *crudRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(new(T)).Error
}
