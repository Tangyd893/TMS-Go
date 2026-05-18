package api

import (
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/health/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, h.service.Health(r.Context()))
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, h.service.Ready(r.Context()))
}
