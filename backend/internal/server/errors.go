package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

func ErrorHandler() fiber.ErrorHandler {
	errorMap := map[apperr.ErrorCode]int{
		apperr.CodeConflict:         http.StatusConflict,
		apperr.CodeValidationFailed: http.StatusBadRequest,
		apperr.CodeInternal:         http.StatusInternalServerError,
		apperr.CodeNotFound:         http.StatusNotFound,
		apperr.CodeUnauthorized:     http.StatusUnauthorized,
		apperr.CodeBadRequest:       http.StatusBadRequest,
		apperr.TooManyRequests:      http.StatusTooManyRequests,
	}

	return func(ctx fiber.Ctx, e error) error {
		status := http.StatusInternalServerError

		if err, ok := errors.AsType[*apperr.InternalError](e); ok {
			if mappedStatus, exists := errorMap[err.Code]; exists {
				status = mappedStatus
			} else {
				slog.Warn("application error code has no mapped http status", "code", err.Code)
			}

			if status >= http.StatusInternalServerError {
				slog.Error("internal server error occurred",
					"err", err.Error(),
					"code", err.Code,
					"details", err.Details,
				)

				return ctx.Status(status).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			response := fiber.Map{
				"message": err.Message,
				"code":    err.Code,
			}

			if len(err.Details) > 0 {
				response["details"] = err.Details
			}

			return ctx.Status(status).JSON(response)
		}

		if err, ok := errors.AsType[*fiber.Error](e); ok {
			return ctx.Status(err.Code).JSON(fiber.Map{"message": err.Message})
		}

		slog.Error("unhandled generic error", "err", e.Error())
		return ctx.Status(status).JSON(fiber.Map{"message": "Internal Server Error"})
	}
}
