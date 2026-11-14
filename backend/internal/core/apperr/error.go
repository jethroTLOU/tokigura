package apperr

import "fmt"

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
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = val
	return e
}
