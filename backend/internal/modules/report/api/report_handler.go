package api

import (
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/report/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	reportService service.ReportService
}

func NewHandler(svc service.ReportService) *Handler {
	return &Handler{reportService: svc}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.reportService.GetDashboard(r.Context())
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, data)
}

func (h *Handler) OrderStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.reportService.GetOrderStats(r.Context())
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, stats)
}

func (h *Handler) TransportEfficiency(w http.ResponseWriter, r *http.Request) {
	eff, err := h.reportService.GetTransportEfficiency(r.Context())
	if err != nil {
		response.Fail(w, r, 500, 500001, err.Error(), nil)
		return
	}
	response.Success(w, r, eff)
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/reports/dashboard", handler.Dashboard)
	mux.HandleFunc("GET /api/v1/reports/order-stats", handler.OrderStats)
	mux.HandleFunc("GET /api/v1/reports/transport-efficiency", handler.TransportEfficiency)
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}
