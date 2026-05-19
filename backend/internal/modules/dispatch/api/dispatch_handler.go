package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	dispatchService service.DispatchService
}

func NewHandler(svc service.DispatchService) *Handler {
	return &Handler{dispatchService: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.dispatchService.List(r.Context(), req)
	if err != nil { response.Fail(w, r, 500, 500001, err.Error(), nil); return }
	response.Success(w, r, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	plan, err := h.dispatchService.GetByID(r.Context(), id)
	if err != nil { response.Fail(w, r, 404, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, plan)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req struct {
		OrderIDs []string `json:"orderIds"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	plan, err := h.dispatchService.Create(r.Context(), req.OrderIDs, userID)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, plan)
}

func (h *Handler) Assign(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var req struct {
		DetailID  string  `json:"detailId"`
		VehicleID *string `json:"vehicleId"`
		DriverID  *string `json:"driverId"`
		CarrierID *string `json:"carrierId"`
	}
	json.Unmarshal(body, &req)
	if err := h.dispatchService.Assign(r.Context(), id, req.DetailID, req.VehicleID, req.DriverID, req.CarrierID); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, nil)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if err := h.dispatchService.Cancel(r.Context(), id); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, nil)
}

func (h *Handler) PendingOrders(w http.ResponseWriter, r *http.Request) {
	result, err := h.dispatchService.GetPendingOrders(r.Context())
	if err != nil { response.Fail(w, r, 500, 500001, err.Error(), nil); return }
	response.Success(w, r, result)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/dispatch/orders", handler.PendingOrders)
	mux.HandleFunc("GET /api/v1/dispatch/plans", handler.List)
	mux.HandleFunc("POST /api/v1/dispatch/plans", handler.Create)
	mux.HandleFunc("GET /api/v1/dispatch/plans/{id}", handler.GetByID)
	mux.HandleFunc("POST /api/v1/dispatch/plans/{id}/assign", handler.Assign)
	mux.HandleFunc("POST /api/v1/dispatch/plans/{id}/cancel", handler.Cancel)
}

func extractID(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' { return path[i+1:] }
	}
	return ""
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok { return e }
	return errors.New(500001, err.Error())
}
