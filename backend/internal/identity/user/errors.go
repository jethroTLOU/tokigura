package user

import "github.com/jethroTLOU/tokigura/internal/core/apperr"

var (
	ErrNotFound      = apperr.New(apperr.CodeNotFound, "user not found")
	ErrAlreadyExists = apperr.New(apperr.CodeConflict, "username or email already exists")
)
