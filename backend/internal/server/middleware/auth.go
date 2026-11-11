package middleware

import (
	"slices"
	"strings"
	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jethroTLOU/tokigura/internal/core/apperr"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
)

type tokenParser interface {
	Parse(tokenString string) (*jwt.Claims, error)
}

type contextKey string

const (
	userIDKey   contextKey = "userId"
	userRoleKey contextKey = "userRole"
)

type AuthMiddleware struct {
	tokens tokenParser
}

func NewAuthMiddleware(tokens *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) Required() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || tokenString == "" {
			return apperr.New(apperr.CodeUnauthorized, "missing or invalid authorization header")
		}

		claims, err := m.tokens.Parse(tokenString)
		if err != nil {
			return apperr.New(apperr.CodeUnauthorized, "invalid or expired token")
		}

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return apperr.New(apperr.CodeUnauthorized, "invalid user id in token payload")
		}

		ctx.Locals(userIDKey, userUUID)
		ctx.Locals(userRoleKey, claims.Role)

		return ctx.Next()
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		role, ok := GetUserRole(ctx)
		if !ok {
			return apperr.New(apperr.CodeUnauthorized, "unauthorized")
		}

		if slices.Contains(allowedRoles, role) {
			return ctx.Next()
		}

		return apperr.New(apperr.CodeUnauthorized, "insufficient permissions")
	}
}

// --- Типобезопасные хелперы для чтения контекста в любых хэндлерах ---

func GetUserID(ctx fiber.Ctx) (uuid.UUID, bool) {
	val, ok := ctx.Locals(userIDKey).(uuid.UUID)
	return val, ok
}

func GetUserRole(ctx fiber.Ctx) (string, bool) {
	val, ok := ctx.Locals(userRoleKey).(string)
	return val, ok
}
