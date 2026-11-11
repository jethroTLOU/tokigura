package jwt

import (
	goJwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	goJwt.RegisteredClaims

	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
