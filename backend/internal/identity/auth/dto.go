package auth

import (
	"net/netip"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/domain"
)

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginDTO struct {
	Username  string
	Password  string
	IP        netip.Addr
	UserAgent string
}

type MeResponse struct {
	ID       uuid.UUID   `json:"id"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
	Email    string      `json:"email"`
	Avatar   *string     `json:"avatar"`
}
