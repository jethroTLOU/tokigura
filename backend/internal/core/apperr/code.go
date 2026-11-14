package apperr

type ErrorCode string

const (
	CodeNotFound         ErrorCode = "not_found"
	CodeValidationFailed ErrorCode = "validation_failed"
	CodeConflict         ErrorCode = "conflict"
	CodeUnauthorized     ErrorCode = "unauthorized"
	CodeForbidden        ErrorCode = "forbidden"
	CodeInternal         ErrorCode = "internal_error"
	CodeBadRequest       ErrorCode = "bad_request"
	CodeTooManyRequests  ErrorCode = "too_many_requests"
)
