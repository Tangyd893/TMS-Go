package dict

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.ListTypes(r.Context())
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, types)
}

func (h *Handler) GetType(w http.ResponseWriter, r *http.Request) {
	code := extractLastPath(r.URL.Path)
	if code == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少字典编码", nil)
		return
	}
	t, err := h.svc.GetType(r.Context(), code)
	if err != nil {
		response.Fail(w, r, http.StatusNotFound, 404001, err.Error(), nil)
		return
	}
	response.Success(w, r, t)
}

func (h *Handler) CreateType(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var t DictType
	if err := json.Unmarshal(body, &t); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}
	result, err := h.svc.CreateType(r.Context(), &t)
	if err != nil {
		appErr := unwrap(err)
		response.Fail(w, r, http.StatusInternalServerError, appErr.Code, appErr.Message, nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request) {
	id := extractLastPath(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var t DictType
	if err := json.Unmarshal(body, &t); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}
	if tID, err := parseUUID(id); err == nil {
		t.ID = tID
	}
	result, err := h.svc.UpdateType(r.Context(), &t)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) DeleteType(w http.ResponseWriter, r *http.Request) {
	id := extractLastPath(r.URL.Path)
	if err := h.svc.DeleteType(r.Context(), id); err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	typeCode := r.URL.Query().Get("typeCode")
	items, err := h.svc.ListItems(r.Context(), typeCode)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, items)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var item DictItem
	if err := json.Unmarshal(body, &item); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}
	result, err := h.svc.CreateItem(r.Context(), &item)
	if err != nil {
		appErr := unwrap(err)
		response.Fail(w, r, http.StatusInternalServerError, appErr.Code, appErr.Message, nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := extractLastPath(r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	var item DictItem
	if err := json.Unmarshal(body, &item); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}
	if tID, err := parseUUID(id); err == nil {
		item.ID = tID
	}
	result, err := h.svc.UpdateItem(r.Context(), &item)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, result)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := extractLastPath(r.URL.Path)
	if err := h.svc.DeleteItem(r.Context(), id); err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, nil)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/dict/types", handler.ListTypes)
	mux.HandleFunc("GET /api/v1/dict/types/{code}", handler.GetType)
	mux.HandleFunc("POST /api/v1/dict/types", handler.CreateType)
	mux.HandleFunc("PUT /api/v1/dict/types/{id}", handler.UpdateType)
	mux.HandleFunc("DELETE /api/v1/dict/types/{id}", handler.DeleteType)
	mux.HandleFunc("GET /api/v1/dict/items", handler.ListItems)
	mux.HandleFunc("POST /api/v1/dict/items", handler.CreateItem)
	mux.HandleFunc("PUT /api/v1/dict/items/{id}", handler.UpdateItem)
	mux.HandleFunc("DELETE /api/v1/dict/items/{id}", handler.DeleteItem)
}

func extractLastPath(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 || idx == len(path)-1 {
		return ""
	}
	return path[idx+1:]
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
