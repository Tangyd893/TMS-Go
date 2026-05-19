package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	financeService service.FinanceService
}

func NewHandler(svc service.FinanceService) *Handler {
	return &Handler{financeService: svc}
}

// -- 应收 --

func (h *Handler) ListReceivables(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	status := r.URL.Query().Get("status")
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.financeService.ListReceivables(r.Context(), req, status)
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

// -- 应付 --

func (h *Handler) ListPayables(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	status := r.URL.Query().Get("status")
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.financeService.ListPayables(r.Context(), req, status)
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

// -- 费用项目 --

func (h *Handler) ListFeeItems(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.financeService.ListFeeItems(r.Context(), req)
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) CreateFeeItem(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var item model.FeeItem
	if err := json.Unmarshal(body, &item); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil)
		return
	}
	result, err := h.financeService.CreateFeeItem(r.Context(), &item)
	if err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

// -- 费用生成 --

type GenerateFeesRequest struct {
	OrderID      string  `json:"orderId"`
	CustomerID   string  `json:"customerId"`
	CarrierID    string  `json:"carrierId"`
	CustomerName string  `json:"customerName"`
	CarrierName  string  `json:"carrierName"`
	Amount       float64 `json:"amount"`
}

func (h *Handler) GenerateFees(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req GenerateFeesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil)
		return
	}
	if err := h.financeService.GenerateFees(r.Context(), req.OrderID, req.CustomerID, req.CarrierID, req.CustomerName, req.CarrierName, req.Amount); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

// -- 对账 --

func (h *Handler) ListStatements(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.financeService.ListStatements(r.Context(), req)
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

type CreateStatementRequest struct {
	PartnerID   string `json:"partnerId"`
	PartnerName string `json:"partnerName"`
	PartnerType string `json:"partnerType"`
	PeriodStart string `json:"periodStart"`
	PeriodEnd   string `json:"periodEnd"`
}

func (h *Handler) CreateStatement(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req CreateStatementRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil)
		return
	}
	result, err := h.financeService.CreateStatement(r.Context(), req.PartnerID, req.PartnerName, req.PartnerType, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) ConfirmStatement(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.financeService.ConfirmStatement(r.Context(), id, userID); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

// -- 结算 --

func (h *Handler) ListSettlements(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.financeService.ListSettlements(r.Context(), req)
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

type CreateSettlementRequest struct {
	StatementID string `json:"statementId"`
}

func (h *Handler) CreateSettlement(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req CreateSettlementRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil)
		return
	}
	result, err := h.financeService.CreateSettlement(r.Context(), req.StatementID)
	if err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) CompleteSettlement(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if err := h.financeService.CompleteSettlement(r.Context(), id); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

// -- 路由注册 --

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/finance/receivables", handler.ListReceivables)
	mux.HandleFunc("GET /api/v1/finance/payables", handler.ListPayables)
	mux.HandleFunc("GET /api/v1/finance/fee-items", handler.ListFeeItems)
	mux.HandleFunc("POST /api/v1/finance/fee-items", handler.CreateFeeItem)
	mux.HandleFunc("POST /api/v1/finance/fees/generate", handler.GenerateFees)
	mux.HandleFunc("GET /api/v1/finance/statements", handler.ListStatements)
	mux.HandleFunc("POST /api/v1/finance/statements", handler.CreateStatement)
	mux.HandleFunc("POST /api/v1/finance/statements/{id}/confirm", handler.ConfirmStatement)
	mux.HandleFunc("GET /api/v1/finance/settlements", handler.ListSettlements)
	mux.HandleFunc("POST /api/v1/finance/settlements", handler.CreateSettlement)
	mux.HandleFunc("POST /api/v1/finance/settlements/{id}/complete", handler.CompleteSettlement)
}

// -- 工具函数 --

func extractID(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return ""
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
