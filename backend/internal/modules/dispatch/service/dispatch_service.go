package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDispatchNotFound = errors.New(404001, "调度计划不存在")
)

type DispatchService interface {
	List(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.DispatchPlan], error)
	GetByID(ctx context.Context, id string) (*model.DispatchPlan, error)
	Create(ctx context.Context, orderIDs []string, userID string) (*model.DispatchPlan, error)
	Assign(ctx context.Context, id string, detailID string, vehicleID, driverID, carrierID *string) error
	Cancel(ctx context.Context, id string) error
	GetPendingOrders(ctx context.Context) (any, error)
}

type dispatchService struct {
	repo repository.DispatchRepository
}

func NewDispatchService(repo repository.DispatchRepository) DispatchService {
	return &dispatchService{repo: repo}
}

func (s *dispatchService) List(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.DispatchPlan], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	plans, total, err := s.repo.List(ctx, req.Page, req.PageSize)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.DispatchPlan]{Items: plans, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *dispatchService) GetByID(ctx context.Context, id string) (*model.DispatchPlan, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, ErrDispatchNotFound }
	plan, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrDispatchNotFound }
		return nil, err
	}
	return plan, nil
}

func (s *dispatchService) Create(ctx context.Context, orderIDs []string, userID string) (*model.DispatchPlan, error) {
	now := time.Now()
	planNo := fmt.Sprintf("DP%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)

	plan := &model.DispatchPlan{
		PlanNo:    planNo,
		Status:    model.DispatchStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if userID != "" {
		uid := uuid.MustParse(userID)
		plan.CreatedBy = &uid
	}

	for _, oid := range orderIDs {
		orderUID := uuid.MustParse(oid)
		plan.Details = append(plan.Details, model.DispatchDetail{
			ID:        uuid.New(),
			OrderID:   orderUID,
			CreatedAt: now,
		})
	}

	if err := s.repo.Create(ctx, plan); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, plan.ID)
}

func (s *dispatchService) Assign(ctx context.Context, id string, detailID string, vehicleID, driverID, carrierID *string) error {
	uid, _ := uuid.Parse(id)
	plan, err := s.repo.FindByID(ctx, uid)
	if err != nil { return ErrDispatchNotFound }

	for i := range plan.Details {
		if plan.Details[i].ID.String() == detailID {
			d := &plan.Details[i]
			if vehicleID != nil {
				vid := uuid.MustParse(*vehicleID)
				d.VehicleID = &vid
			}
			if driverID != nil {
				did := uuid.MustParse(*driverID)
				d.DriverID = &did
			}
			if carrierID != nil {
				cid := uuid.MustParse(*carrierID)
				d.CarrierID = &cid
			}
			if err := s.repo.UpdateDetail(ctx, d.ID, d); err != nil {
				return err
			}
		}
	}

	plan.Status = model.DispatchStatusAssigned
	plan.UpdatedAt = time.Now()

	if err := s.repo.UpdateStatus(ctx, uid, model.DispatchStatusAssigned); err != nil {
		return err
	}
	return nil
}

func (s *dispatchService) Cancel(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)
	return s.repo.UpdateStatus(ctx, uid, model.DispatchStatusCancelled)
}

func (s *dispatchService) GetPendingOrders(ctx context.Context) (any, error) {
	return s.repo.FindPendingOrders(ctx)
}
