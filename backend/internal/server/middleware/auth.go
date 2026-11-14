package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
	"github.com/jethroTLOU/tokigura/internal/core/domain"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
)

type tokenParser interface {
	Parse(tokenString string) (*jwt.Claims, error)
}

type (
	userIDCtxKey   struct{}
	userRoleCtxKey struct{}
)

type AuthMiddleware struct {
	tokens tokenParser
}

func NewAuthMiddleware(tokens *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || tokenString == "" {
			writeAuthError(
				w,
				http.StatusUnauthorized,
				apperr.CodeUnauthorized,
				"missing or invalid authorization header",
			)
			return
		}

		claims, err := m.tokens.Parse(tokenString)
		if err != nil {
			writeAuthError(
				w,
				http.StatusUnauthorized,
				apperr.CodeUnauthorized,
				"invalid or expired token",
			)
			return
		}

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			writeAuthError(
				w,
				http.StatusUnauthorized,
				apperr.CodeUnauthorized,
				"invalid user id in token payload",
			)
			return
		}

		ctx := context.WithValue(r.Context(), userIDCtxKey{}, userUUID)
		ctx = context.WithValue(ctx, userRoleCtxKey{}, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetUserRole(r.Context())
			if !ok {
				writeAuthError(w, http.StatusUnauthorized, apperr.CodeUnauthorized, "unauthorized")
				return
			}
			if !slices.Contains(allowedRoles, role) {
				writeAuthError(
					w,
					http.StatusForbidden,
					apperr.CodeForbidden,
					"insufficient permissions",
				)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(userIDCtxKey{}).(uuid.UUID)
	return val, ok
}

func GetUserRole(ctx context.Context) (domain.Role, bool) {
	val, ok := ctx.Value(userRoleCtxKey{}).(domain.Role)
	return val, ok
}

func writeAuthError(w http.ResponseWriter, status int, code apperr.ErrorCode, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": msg,
		"code":    code,
	})
}
