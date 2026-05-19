package dict

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListTypes(ctx context.Context) ([]DictType, error) {
	var types []DictType
	err := r.db.WithContext(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_no asc")
	}).Order("created_at desc").Find(&types).Error
	return types, err
}

func (r *Repository) FindTypeByCode(ctx context.Context, code string) (*DictType, error) {
	var t DictType
	err := r.db.WithContext(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_no asc")
	}).Where("code = ?", code).First(&t).Error
	return &t, err
}

func (r *Repository) CreateType(ctx context.Context, t *DictType) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Repository) UpdateType(ctx context.Context, t *DictType) error {
	return r.db.WithContext(ctx).Model(t).Where("id = ?", t.ID).Updates(map[string]any{
		"name": t.Name, "status": t.Status, "remark": t.Remark, "updated_at": gorm.Expr("now()"),
	}).Error
}

func (r *Repository) DeleteType(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&DictType{}).Error
}

func (r *Repository) ListItems(ctx context.Context, typeCode string) ([]DictItem, error) {
	var items []DictItem
	err := r.db.WithContext(ctx).Where("type_code = ?", typeCode).Order("sort_no asc").Find(&items).Error
	return items, err
}

func (r *Repository) CreateItem(ctx context.Context, item *DictItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) UpdateItem(ctx context.Context, item *DictItem) error {
	return r.db.WithContext(ctx).Model(item).Where("id = ?", item.ID).Updates(map[string]any{
		"item_name": item.ItemName, "sort_no": item.SortNo, "status": item.Status, "remark": item.Remark, "updated_at": gorm.Expr("now()"),
	}).Error
}

func (r *Repository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&DictItem{}).Error
}
