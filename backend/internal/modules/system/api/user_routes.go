package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/v1/users", handler.List)
	mux.HandleFunc("GET /api/v1/users/{id}", handler.GetByID)
	mux.HandleFunc("POST /api/v1/users", handler.Create)
	mux.HandleFunc("PUT /api/v1/users/{id}", handler.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", handler.Delete)
}
