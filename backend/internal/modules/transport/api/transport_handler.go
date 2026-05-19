package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	transportService service.TransportService
}

func NewHandler(svc service.TransportService) *Handler {
	return &Handler{transportService: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.transportService.List(r.Context(), req, keyword, status)
	if err != nil { response.Fail(w, r, 500, 500001, err.Error(), nil); return }
	response.Success(w, r, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	task, err := h.transportService.GetByID(r.Context(), id)
	if err != nil { response.Fail(w, r, 404, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, task)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var task model.TransportTask
	if err := json.Unmarshal(body, &task); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	result, err := h.transportService.Create(r.Context(), &task)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, result)
}

func (h *Handler) Depart(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	task, err := h.transportService.Depart(r.Context(), id)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, task)
}

func (h *Handler) Arrive(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	task, err := h.transportService.Arrive(r.Context(), id)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, task)
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	task, err := h.transportService.Sign(r.Context(), id)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, task)
}

func (h *Handler) AddNode(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var req model.CreateNodeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	node, err := h.transportService.AddNode(r.Context(), id, req, userID)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, node)
}

func (h *Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	nodes, err := h.transportService.GetNodes(r.Context(), id)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, nodes)
}

func (h *Handler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var req model.CreateReceiptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	receipt, err := h.transportService.CreateReceipt(r.Context(), id, req, userID)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, receipt)
}

func (h *Handler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	receipt, err := h.transportService.GetReceipt(r.Context(), id)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, receipt)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/transport/tasks", handler.List)
	mux.HandleFunc("POST /api/v1/transport/tasks", handler.Create)
	mux.HandleFunc("GET /api/v1/transport/tasks/{id}", handler.GetByID)
	mux.HandleFunc("POST /api/v1/transport/tasks/{id}/depart", handler.Depart)
	mux.HandleFunc("POST /api/v1/transport/tasks/{id}/arrive", handler.Arrive)
	mux.HandleFunc("POST /api/v1/transport/tasks/{id}/sign", handler.Sign)
	mux.HandleFunc("POST /api/v1/transport/tasks/{id}/nodes", handler.AddNode)
	mux.HandleFunc("GET /api/v1/transport/tasks/{id}/nodes", handler.GetNodes)
	mux.HandleFunc("POST /api/v1/transport/tasks/{id}/receipt", handler.CreateReceipt)
	mux.HandleFunc("GET /api/v1/transport/tasks/{id}/receipt", handler.GetReceipt)
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
