package response

import (
	"encoding/json"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/platform/requestid"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"traceId,omitempty"`
}

func Success(w http.ResponseWriter, r *http.Request, data any) {
	writeJSON(w, http.StatusOK, Body{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: requestid.FromContext(r.Context()),
	})
}

func Fail(w http.ResponseWriter, r *http.Request, status int, code int, message string, data any) {
	writeJSON(w, status, Body{
		Code:    code,
		Message: message,
		Data:    data,
		TraceID: requestid.FromContext(r.Context()),
	})
}

func writeJSON(w http.ResponseWriter, status int, body Body) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
