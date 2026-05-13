package apperr

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type InternalError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Err     error
}

func NewValidation(details map[string]any) *InternalError {
	return &InternalError{
		Code:    CodeValidationFailed,
		Message: "Validation error.",
		Details: details,
	}
}

func FromValidator(e error) *InternalError {
	err, ok := errors.AsType[*validator.ValidationErrors](e)
	if !ok {
		return Wrap(e, CodeInternal, "Given non validation error type")
	}

	errs := make(map[string]any, len(*err))

	for _, validationError := range *err {
		errs[validationError.Field()] = validationError.Error()
	}

	return NewValidation(errs)
}

func NewBindError(err error) *InternalError {
	return Wrap(err, CodeBadRequest, "Invalid request format. Please check the submitted data.")
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string) *InternalError {
	return &InternalError{
		Code:    code,
		Message: message,
	}
}

func (e *InternalError) With(key, val string) *InternalError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = val
	return e
}

func Wrap(err error, code ErrorCode, message string) *InternalError {
	return &InternalError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
