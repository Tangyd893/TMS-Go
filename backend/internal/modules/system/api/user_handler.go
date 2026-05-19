package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	keyword := r.URL.Query().Get("keyword")

	req := pagination.PageRequest{Page: page, PageSize: pageSize}
	result, err := h.userService.List(r.Context(), req, keyword)
	if err != nil {
		response.Fail(w, r, http.StatusInternalServerError, 500001, err.Error(), nil)
		return
	}

	response.Success(w, r, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r.URL.Path, "/api/v1/users/")
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少用户ID", nil)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusNotFound, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, user)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var params model.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	if params.Username == "" || params.Password == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "用户名和密码不能为空", nil)
		return
	}

	user, err := h.userService.Create(r.Context(), params)
	if err != nil {
		appErr := unwrapAppError(err)
		status := http.StatusInternalServerError
		if appErr.Code == 409001 {
			status = http.StatusConflict
		}
		response.Fail(w, r, status, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r.URL.Path, "/api/v1/users/")
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少用户ID", nil)
		return
	}

	var params model.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	user, err := h.userService.Update(r.Context(), id, params)
	if err != nil {
		appErr := unwrapAppError(err)
		status := http.StatusInternalServerError
		if appErr.Code == 404001 {
			status = http.StatusNotFound
		} else if appErr.Code == 409001 {
			status = http.StatusConflict
		}
		response.Fail(w, r, status, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r.URL.Path, "/api/v1/users/")
	if id == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "缺少用户ID", nil)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusNotFound, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, nil)
}

func extractPathParam(path, prefix string) string {
	trimmed := strings.TrimPrefix(path, prefix)
	slashIdx := strings.Index(trimmed, "/")
	if slashIdx >= 0 {
		return trimmed[:slashIdx]
	}
	return trimmed
}

func unwrapAppError(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
