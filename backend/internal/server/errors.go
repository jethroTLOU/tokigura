package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Wrap(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			HandleError(w, r, err)
		}
	}
}

var errorMap = map[apperr.ErrorCode]int{
	apperr.CodeConflict:         http.StatusConflict,
	apperr.CodeValidationFailed: http.StatusBadRequest,
	apperr.CodeInternal:         http.StatusInternalServerError,
	apperr.CodeNotFound:         http.StatusNotFound,
	apperr.CodeUnauthorized:     http.StatusUnauthorized,
	apperr.CodeForbidden:        http.StatusForbidden,
	apperr.CodeBadRequest:       http.StatusBadRequest,
	apperr.CodeTooManyRequests:  http.StatusTooManyRequests,
}

func HandleError(w http.ResponseWriter, _ *http.Request, e error) {
	status := http.StatusInternalServerError

	if err, ok := errors.AsType[*apperr.AppError](e); ok {
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
			_ = WriteJSON(w, status, Map{"message": "Internal Server Error"})
			return
		}

		response := Map{
			"message": err.Message,
			"code":    err.Code,
		}
		if len(err.Details) > 0 {
			response["details"] = err.Details
		}

		_ = WriteJSON(w, status, response)
		return
	}

	slog.Error("unhandled generic error", "err", e.Error())
	_ = WriteJSON(w, status, Map{"message": "Internal Server Error"})
}
