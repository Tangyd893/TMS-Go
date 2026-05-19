package repository

import (
	"context"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.Order, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateStatusHistory(ctx context.Context, h *model.StatusHistory) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Order{}).Preload("CargoItems")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("order_no LIKE ? OR customer_name LIKE ? OR origin_name LIKE ? OR dest_name LIKE ?", like, like, like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Preload("CargoItems").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(order).Omit("CargoItems").Updates(order).Error; err != nil {
			return err
		}
		if len(order.CargoItems) > 0 {
			tx.Where("order_id = ?", order.ID).Delete(&model.OrderCargo{})
			for i := range order.CargoItems {
				order.CargoItems[i].ID = uuid.New()
				order.CargoItems[i].OrderID = order.ID
				order.CargoItems[i].CreatedAt = time.Now()
			}
			return tx.Create(&order.CargoItems).Error
		}
		return nil
	})
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Order{}).Error
}

func (r *orderRepository) CreateStatusHistory(ctx context.Context, h *model.StatusHistory) error {
	return r.db.WithContext(ctx).Create(h).Error
}
