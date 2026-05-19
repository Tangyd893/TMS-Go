package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/platform/requestid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestid.Header)
		if id == "" {
			id = newRequestID()
		}

		w.Header().Set(requestid.Header, id)
		ctx := requestid.WithContext(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "request-id-unavailable"
	}
	return hex.EncodeToString(bytes[:])
}
