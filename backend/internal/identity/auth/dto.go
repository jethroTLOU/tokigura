package auth

import (
	"net/netip"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/domain"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
