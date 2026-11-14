package apperr

import (
	"errors"
	"fmt"
	"maps"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Err     error
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewValidation(details map[string]any) *AppError {
	return &AppError{
		Code:    CodeValidationFailed,
		Message: "Validation error.",
		Details: details,
	}
}

func NewBindError(err error) *AppError {
	return Wrap(err, CodeBadRequest, "Invalid request format. Please check the submitted data.")
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) With(key string, val any) *AppError {
	details := maps.Clone(e.Details)
	if details == nil {
		details = make(map[string]any, 1)
	}

	details[key] = val

	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
		Err:     e.Err,
	}
}

func IsCode(err error, code ErrorCode) bool {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr.Code == code
	}
	return false
}
