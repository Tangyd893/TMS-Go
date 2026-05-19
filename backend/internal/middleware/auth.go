package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type userIDKey struct{}

func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Fail(w, r, http.StatusUnauthorized, 401001, "未登录或令牌无效", nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Fail(w, r, http.StatusUnauthorized, 401001, "未登录或令牌无效", nil)
				return
			}

			claims, err := jwtManager.ParseAccessToken(parts[1])
			if err != nil {
				response.Fail(w, r, http.StatusUnauthorized, 401003, "令牌已过期", nil)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey{}, claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) string {
	value, ok := ctx.Value(userIDKey{}).(string)
	if !ok {
		return ""
	}
	return value
}
