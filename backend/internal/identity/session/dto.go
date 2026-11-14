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

type RotateSessionDTO struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	Role      user.Role
	OldHash   string
	IP        netip.Addr
}
type UpdateRotatedParams struct {
	ID        uuid.UUID
	OldHash   string
	NewHash   string
	IP        netip.Addr
	ExpiresAt time.Time
}
