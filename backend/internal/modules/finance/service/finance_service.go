package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
)

var (
	ErrStatementNotFound = errors.New(404001, "对账单不存在")
)

type FinanceService interface {
	// 应收
	ListReceivables(ctx context.Context, req pagination.PageRequest, status string) (*pagination.PageResult[model.Receivable], error)
	// 应付
	ListPayables(ctx context.Context, req pagination.PageRequest, status string) (*pagination.PageResult[model.Payable], error)
	// 费用项目
	ListFeeItems(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.FeeItem], error)
	CreateFeeItem(ctx context.Context, item *model.FeeItem) (*model.FeeItem, error)
	// 费用生成
	GenerateFees(ctx context.Context, orderID, customerID, carrierID string, customerName, carrierName string, amount float64) error
	// 对账
	ListStatements(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.Statement], error)
	CreateStatement(ctx context.Context, partnerID, partnerName, partnerType, periodStart, periodEnd string) (*model.Statement, error)
	ConfirmStatement(ctx context.Context, id, userID string) error
	// 结算
	ListSettlements(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.Settlement], error)
	CreateSettlement(ctx context.Context, statementID string) (*model.Settlement, error)
	CompleteSettlement(ctx context.Context, id string) error
}

type financeService struct {
	repo repository.FinanceRepository
}

func NewFinanceService(repo repository.FinanceRepository) FinanceService {
	return &financeService{repo: repo}
}

func (s *financeService) ListReceivables(ctx context.Context, req pagination.PageRequest, status string) (*pagination.PageResult[model.Receivable], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.ListReceivables(ctx, req.Page, req.PageSize, status)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.Receivable]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *financeService) ListPayables(ctx context.Context, req pagination.PageRequest, status string) (*pagination.PageResult[model.Payable], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.ListPayables(ctx, req.Page, req.PageSize, status)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.Payable]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *financeService) ListFeeItems(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.FeeItem], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.ListFeeItems(ctx, req.Page, req.PageSize)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.FeeItem]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *financeService) CreateFeeItem(ctx context.Context, item *model.FeeItem) (*model.FeeItem, error) {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	if err := s.repo.CreateFeeItem(ctx, item); err != nil { return nil, err }
	return item, nil
}

func (s *financeService) GenerateFees(ctx context.Context, orderID, customerID, carrierID string, customerName, carrierName string, amount float64) error {
	orderUID := uuid.MustParse(orderID)
	now := time.Now()

	if customerID != "" {
		custUID := uuid.MustParse(customerID)
		receivableNo := fmt.Sprintf("AR%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
		rv := &model.Receivable{
			ReceivableNo: receivableNo,
			OrderID:      orderUID,
			CustomerID:   custUID,
			CustomerName: customerName,
			FeeItemName:  "运费",
			Amount:       amount,
			TotalAmount:  amount,
			Status:       model.ReceivableStatusPending,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.repo.GenerateReceivable(ctx, rv); err != nil { return err }
	}

	if carrierID != "" {
		carrierUID := uuid.MustParse(carrierID)
		payableNo := fmt.Sprintf("AP%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
		pv := &model.Payable{
			PayableNo:   payableNo,
			OrderID:     orderUID,
			CarrierID:   carrierUID,
			CarrierName: carrierName,
			FeeItemName: "运费",
			Amount:      amount * 0.8,
			TotalAmount: amount * 0.8,
			Status:      string(ReceivableStatusSettled),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := s.repo.GeneratePayable(ctx, pv); err != nil { return err }
	}

	return nil
}

func (s *financeService) ListStatements(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.Statement], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.ListStatements(ctx, req.Page, req.PageSize)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.Statement]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *financeService) CreateStatement(ctx context.Context, partnerID, partnerName, partnerType, periodStart, periodEnd string) (*model.Statement, error) {
	now := time.Now()
	stmtNo := fmt.Sprintf("ST%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)

	partnerUID := uuid.MustParse(partnerID)
	stmt := &model.Statement{
		StatementNo:          stmtNo,
		PartnerID:            partnerUID,
		PartnerName:          partnerName,
		PartnerType:          partnerType,
		StatementPeriodStart: periodStart,
		StatementPeriodEnd:   periodEnd,
		Status:               model.StatementStatusDraft,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := s.repo.CreateStatement(ctx, stmt); err != nil { return nil, err }
	return stmt, nil
}

func (s *financeService) ConfirmStatement(ctx context.Context, id, userID string) error {
	uid := uuid.MustParse(id)
	userUID := uuid.MustParse(userID)
	return s.repo.ConfirmStatement(ctx, uid, userUID)
}

func (s *financeService) ListSettlements(ctx context.Context, req pagination.PageRequest) (*pagination.PageResult[model.Settlement], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.ListSettlements(ctx, req.Page, req.PageSize)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.Settlement]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *financeService) CreateSettlement(ctx context.Context, statementID string) (*model.Settlement, error) {
	now := time.Now()
	settleNo := fmt.Sprintf("STL%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
	stmtUID := uuid.MustParse(statementID)

	settle := &model.Settlement{
		SettlementNo:  settleNo,
		StatementID:   stmtUID,
		PartnerID:     uuid.New(),
		PartnerType:   "customer",
		Status:        model.SettlementStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.CreateSettlement(ctx, settle); err != nil { return nil, err }
	return settle, nil
}

func (s *financeService) CompleteSettlement(ctx context.Context, id string) error {
	uid := uuid.MustParse(id)
	return s.repo.CompleteSettlement(ctx, uid)
}

// 使用自身常量
const ReceivableStatusSettled = "settled"
