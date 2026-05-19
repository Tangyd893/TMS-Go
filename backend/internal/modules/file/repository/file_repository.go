package repository

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/file/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, f *model.FileObject) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.FileObject, error)
	ListByBiz(ctx context.Context, bizType string, bizID uuid.UUID) ([]model.FileObject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, f *model.FileObject) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.FileObject, error) {
	var f model.FileObject
	err := r.db.WithContext(ctx).First(&f, "id = ?", id).Error
	return &f, err
}

func (r *fileRepository) ListByBiz(ctx context.Context, bizType string, bizID uuid.UUID) ([]model.FileObject, error) {
	var files []model.FileObject
	err := r.db.WithContext(ctx).
		Where("biz_type = ? AND biz_id = ?", bizType, bizID).
		Order("created_at desc").
		Find(&files).Error
	return files, err
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.FileObject{}, "id = ?", id).Error
}
