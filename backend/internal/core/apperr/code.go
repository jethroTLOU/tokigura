package apperr

type ErrorCode string

const (
	CodeNotFound         ErrorCode = "not_found"
	CodeValidationFailed ErrorCode = "validation_failed"
	CodeConflict         ErrorCode = "conflict"
	CodeUnauthorized     ErrorCode = "unauthorized"
	CodeInternal         ErrorCode = "internal_error"
)
