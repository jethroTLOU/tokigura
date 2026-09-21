package auth

import (
	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

func (h *Handler) AuthRequired() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := string(ctx.Request().Header.Peek("Authorization"))
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return apperr.New(apperr.CodeUnauthorized, "missing or invalid token")
		}

		tokenString := authHeader[7:]

		claims, err := h.tokenParser.Parse(tokenString)
		if err != nil {
			return apperr.New(apperr.CodeUnauthorized, "invalid token")
		}

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return apperr.New(apperr.CodeUnauthorized, "invalid token payload: corrupted user_id")
		}

		ctx.Locals("userId", userUUID)

		return ctx.Next()
	}
}
