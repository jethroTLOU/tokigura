package session

import (
	"net/netip"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/identity/user"
)

type CreateSessionDTO struct {
	UserID    uuid.UUID
	Role      user.Role
	IP        netip.Addr
	UserAgent string
}

type UpdateRotatedParams struct {
	ID        uuid.UUID
	Hash      string
	IP        netip.Addr
	ExpiresAt time.Time
}
