package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	exceptionService service.ExceptionService
}

func NewHandler(svc service.ExceptionService) *Handler {
	return &Handler{exceptionService: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	exceptionType := r.URL.Query().Get("exceptionType")
	status := r.URL.Query().Get("status")
	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.exceptionService.List(r.Context(), req, exceptionType, status)
	if err != nil { response.Fail(w, r, 500, 500001, err.Error(), nil); return }
	response.Success(w, r, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	e, err := h.exceptionService.GetByID(r.Context(), id)
	if err != nil { response.Fail(w, r, 404, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, e)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req model.CreateExceptionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	e, err := h.exceptionService.Create(r.Context(), req, userID, "")
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, e)
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var req model.HandleExceptionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.Fail(w, r, 400, 400001, "请求参数错误", nil); return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.exceptionService.Handle(r.Context(), id, userID, "", req.HandleResult); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, nil)
}

func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if err := h.exceptionService.Close(r.Context(), id); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, nil)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// 路径格式: /api/v1/transport/tasks/{id}/exceptions/{id}/logs
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		response.Fail(w, r, 400, 400001, "缺少异常ID", nil); return
	}
	exceptionID := parts[len(parts)-2]
	logs, err := h.exceptionService.GetLogs(r.Context(), exceptionID)
	if err != nil { response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return }
	response.Success(w, r, logs)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/transport/exceptions", handler.List)
	mux.HandleFunc("POST /api/v1/transport/exceptions", handler.Create)
	mux.HandleFunc("GET /api/v1/transport/exceptions/{id}", handler.GetByID)
	mux.HandleFunc("PUT /api/v1/transport/exceptions/{id}/handle", handler.Handle)
	mux.HandleFunc("POST /api/v1/transport/exceptions/{id}/close", handler.Close)
	mux.HandleFunc("GET /api/v1/transport/exceptions/{id}/logs", handler.GetLogs)
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
