package auth

import "github.com/jethroTLOU/tokigura/internal/core/apperr"

var (
	ErrInvalidCredentials = apperr.New(apperr.CodeUnauthorized, "invalid username or password")
	ErrSessionExpired     = apperr.New(apperr.CodeUnauthorized, "session expired")
	ErrUserNotFound       = apperr.New(apperr.CodeUnauthorized, "user not found")
)
