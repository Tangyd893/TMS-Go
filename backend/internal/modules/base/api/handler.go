package api

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/base/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type CrudHandler[T any] struct {
	svc service.EntityService[T]
}

func NewCrudHandler[T any](svc service.EntityService[T]) *CrudHandler[T] {
	return &CrudHandler[T]{svc: svc}
}

func (h *CrudHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	keyword := r.URL.Query().Get("keyword")

	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.svc.List(r.Context(), req, keyword)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *CrudHandler[T]) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少ID", nil)
		return
	}
	entity, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Fail(w, r, http.StatusNotFound, 404001, err.Error(), nil)
		return
	}
	response.Success(w, r, entity)
}

func (h *CrudHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "读取请求体失败", nil)
		return
	}
	entity, err := h.svc.Create(r.Context(), body)
	if err != nil {
		appErr := unwrapError(err)
		response.Fail(w, r, http.StatusInternalServerError, appErr.Code, appErr.Message, nil)
		return
	}
	response.Success(w, r, entity)
}

func (h *CrudHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少ID", nil)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "读取请求体失败", nil)
		return
	}
	entity, err := h.svc.Update(r.Context(), id, body)
	if err != nil {
		appErr := unwrapError(err)
		response.Fail(w, r, http.StatusInternalServerError, appErr.Code, appErr.Message, nil)
		return
	}
	response.Success(w, r, entity)
}

func (h *CrudHandler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少ID", nil)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

func RegisterCrudRoutes[T any](mux *http.ServeMux, prefix string, handler *CrudHandler[T]) {
	mux.HandleFunc("GET "+prefix, handler.List)
	mux.HandleFunc("POST "+prefix, handler.Create)
	mux.HandleFunc("GET "+prefix+"/{id}", handler.GetByID)
	mux.HandleFunc("PUT "+prefix+"/{id}", handler.Update)
	mux.HandleFunc("DELETE "+prefix+"/{id}", handler.Delete)
}

func extractID(path string) string {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash < 0 {
		return ""
	}
	return path[lastSlash+1:]
}

func unwrapError(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
