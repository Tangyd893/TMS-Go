package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready)
	mux.HandleFunc("GET /api/v1/health", handler.Health)
	mux.HandleFunc("GET /api/v1/ready", handler.Ready)
}
