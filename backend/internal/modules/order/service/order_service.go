package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/domain"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New(404001, "订单不存在")
)

type OrderService interface {
	List(ctx context.Context, req pagination.PageRequest, keyword, status string) (*pagination.PageResult[model.Order], error)
	GetByID(ctx context.Context, id string) (*model.Order, error)
	Create(ctx context.Context, req model.CreateOrderRequest, userID string) (*model.Order, error)
	Update(ctx context.Context, id string, req model.UpdateOrderRequest) (*model.Order, error)
	Submit(ctx context.Context, id, userID string) (*model.Order, error)
	Cancel(ctx context.Context, id, userID string) (*model.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) List(ctx context.Context, req pagination.PageRequest, keyword, status string) (*pagination.PageResult[model.Order], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }

	orders, total, err := s.repo.List(ctx, req.Page, req.PageSize, keyword, status)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[model.Order]{
		Items: orders, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*model.Order, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, ErrOrderNotFound }
	order, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrOrderNotFound }
		return nil, err
	}
	return order, nil
}

func (s *orderService) Create(ctx context.Context, req model.CreateOrderRequest, userID string) (*model.Order, error) {
	now := time.Now()
	orderNo := fmt.Sprintf("TMS%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)

	order := &model.Order{
		OrderNo:           orderNo,
		CustomerID:        req.CustomerID,
		ShipperName:       req.ShipperName,
		ShipperPhone:      req.ShipperPhone,
		ShipperAddress:    req.ShipperAddress,
		ReceiverName:      req.ReceiverName,
		ReceiverPhone:     req.ReceiverPhone,
		ReceiverAddress:   req.ReceiverAddress,
		OriginName:        req.OriginName,
		DestName:          req.DestName,
		PlanPickupTime:    req.PlanPickupTime,
		PlanDeliveryTime:  req.PlanDeliveryTime,
		TransportRequirement: req.TransportRequirement,
		Remark:            req.Remark,
		Status:            domain.OrderStatusDraft,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if userID != "" {
		uid := uuid.MustParse(userID)
		order.CreatedBy = &uid
	}

	for _, ci := range req.CargoItems {
		cargo := model.OrderCargo{
			ID:        uuid.New(),
			CargoName: ci.CargoName,
			CargoType: ci.CargoType,
			Quantity:  ci.Quantity,
			Weight:    ci.Weight,
			Volume:    ci.Volume,
			Unit:      ci.Unit,
			CreatedAt: now,
		}
		if cargo.Quantity == 0 { cargo.Quantity = 1 }
		order.CargoItems = append(order.CargoItems, cargo)
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, order.ID)
}

func (s *orderService) Update(ctx context.Context, id string, req model.UpdateOrderRequest) (*model.Order, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, ErrOrderNotFound }

	order, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrOrderNotFound }

	if order.Status != domain.OrderStatusDraft {
		return nil, errors.New(610001, "仅草稿订单可编辑")
	}

	if req.ShipperName != nil { order.ShipperName = *req.ShipperName }
	if req.ShipperPhone != nil { order.ShipperPhone = *req.ShipperPhone }
	if req.ShipperAddress != nil { order.ShipperAddress = *req.ShipperAddress }
	if req.ReceiverName != nil { order.ReceiverName = *req.ReceiverName }
	if req.ReceiverPhone != nil { order.ReceiverPhone = *req.ReceiverPhone }
	if req.ReceiverAddress != nil { order.ReceiverAddress = *req.ReceiverAddress }
	if req.OriginName != nil { order.OriginName = *req.OriginName }
	if req.DestName != nil { order.DestName = *req.DestName }
	if req.PlanPickupTime != nil { order.PlanPickupTime = req.PlanPickupTime }
	if req.PlanDeliveryTime != nil { order.PlanDeliveryTime = req.PlanDeliveryTime }
	if req.TransportRequirement != nil { order.TransportRequirement = *req.TransportRequirement }
	if req.Remark != nil { order.Remark = *req.Remark }
	if req.CargoItems != nil {
		order.CargoItems = nil
		for _, ci := range req.CargoItems {
			c := model.OrderCargo{ID: uuid.New(), OrderID: order.ID, CargoName: ci.CargoName, CargoType: ci.CargoType, Quantity: ci.Quantity, Weight: ci.Weight, Volume: ci.Volume, Unit: ci.Unit, CreatedAt: time.Now()}
			if c.Quantity == 0 { c.Quantity = 1 }
			order.CargoItems = append(order.CargoItems, c)
		}
	}

	if err := s.repo.Update(ctx, order); err != nil { return nil, err }
	return s.repo.FindByID(ctx, uid)
}

func (s *orderService) Submit(ctx context.Context, id, userID string) (*model.Order, error) {
	uid, _ := uuid.Parse(id)
	order, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrOrderNotFound }

	sm := domain.NewOrderStateMachine(order.Status)
	if err := sm.CanSubmit(); err != nil { return nil, err }

	fromStatus := order.Status
	newStatus := sm.Submit()
	order.Status = newStatus

	if err := s.repo.UpdateStatus(ctx, order.ID, newStatus); err != nil { return nil, err }

	uidOp := uuid.MustParse(userID)
	s.repo.CreateStatusHistory(ctx, &model.StatusHistory{
		BizType: "order", BizID: order.ID, FromStatus: fromStatus, ToStatus: newStatus,
		Action: "submit", OperatorID: &uidOp,
	})

	return order, nil
}

func (s *orderService) Cancel(ctx context.Context, id, userID string) (*model.Order, error) {
	uid, _ := uuid.Parse(id)
	order, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrOrderNotFound }

	sm := domain.NewOrderStateMachine(order.Status)
	if err := sm.CanCancel(); err != nil { return nil, err }

	fromStatus := order.Status
	newStatus := sm.Cancel()
	order.Status = newStatus

	if err := s.repo.UpdateStatus(ctx, order.ID, newStatus); err != nil { return nil, err }

	uidOp := uuid.MustParse(userID)
	s.repo.CreateStatusHistory(ctx, &model.StatusHistory{
		BizType: "order", BizID: order.ID, FromStatus: fromStatus, ToStatus: newStatus,
		Action: "cancel", OperatorID: &uidOp,
	})

	return order, nil
}
