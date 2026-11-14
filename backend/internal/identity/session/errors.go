package session

import "github.com/jethroTLOU/tokigura/internal/core/apperr"

var ErrNotFound = apperr.New(apperr.CodeUnauthorized, "invalid session")
