package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	orderService service.OrderService
}

func NewHandler(svc service.OrderService) *Handler {
	return &Handler{orderService: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")

	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.orderService.List(r.Context(), req, keyword, status)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, http.StatusBadRequest, 400001, "缺少订单ID", nil); return }
	order, err := h.orderService.GetByID(r.Context(), id)
	if err != nil {
		response.Fail(w, r, http.StatusNotFound, unwrap(err).Code, err.Error(), nil)
		return
	}
	response.Success(w, r, order)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req model.CreateOrderRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	order, err := h.orderService.Create(r.Context(), req, userID)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, order)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, http.StatusBadRequest, 400001, "缺少订单ID", nil); return }
	body, _ := io.ReadAll(r.Body)
	var req model.UpdateOrderRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil); return
	}
	order, err := h.orderService.Update(r.Context(), id, req)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, order)
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, http.StatusBadRequest, 400001, "缺少订单ID", nil); return }
	userID := middleware.UserIDFromContext(r.Context())
	order, err := h.orderService.Submit(r.Context(), id, userID)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, order)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, http.StatusBadRequest, 400001, "缺少订单ID", nil); return }
	userID := middleware.UserIDFromContext(r.Context())
	order, err := h.orderService.Cancel(r.Context(), id, userID)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, order)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/orders", handler.List)
	mux.HandleFunc("POST /api/v1/orders", handler.Create)
	mux.HandleFunc("GET /api/v1/orders/{id}", handler.GetByID)
	mux.HandleFunc("PUT /api/v1/orders/{id}", handler.Update)
	mux.HandleFunc("POST /api/v1/orders/{id}/submit", handler.Submit)
	mux.HandleFunc("POST /api/v1/orders/{id}/cancel", handler.Cancel)
}

func extractID(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 || idx == len(path)-1 { return "" }
	return path[idx+1:]
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok { return e }
	return errors.New(500001, err.Error())
}
