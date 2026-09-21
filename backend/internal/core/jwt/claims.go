package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
)

type Claims struct {
	jwt.RegisteredClaims

	UserID string    `json:"user_id"`
	Role   user.Role `json:"role"`
}
