package jwt

import (
	goJwt "github.com/golang-jwt/jwt/v5"
	"github.com/jethroTLOU/tokigura/internal/core/domain"
)

type Claims struct {
	goJwt.RegisteredClaims

	UserID string      `json:"user_id"`
	Role   domain.Role `json:"role"`
}
