package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	healthapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/api"
	healthsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/service"
)

func NewHTTPServer(cfg config.Config, log *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	healthService := healthsvc.NewService()
	healthHandler := healthapi.NewHandler(healthService)
	healthapi.RegisterRoutes(mux, healthHandler)

	handler := middleware.RequestID(middleware.AccessLog(log)(mux))

	return &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: handler,
	}
}
