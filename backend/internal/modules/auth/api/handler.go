package api

import (
	"encoding/json"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	if req.Username == "" || req.Password == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "用户名和密码不能为空", nil)
		return
	}

	result, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusUnauthorized, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, result)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusUnauthorized, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, result)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		response.Fail(w, r, http.StatusUnauthorized, 401001, "未登录或令牌无效", nil)
		return
	}

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.Fail(w, r, http.StatusUnauthorized, 401004, "获取用户信息失败", nil)
		return
	}

	response.Success(w, r, user)
}

func unwrapAppError(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
