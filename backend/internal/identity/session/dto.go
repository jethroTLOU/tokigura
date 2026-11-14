package session

import (
	"net/netip"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/domain"
)

type CreateSessionDTO struct {
	UserID    uuid.UUID
	Role      domain.Role
	IP        netip.Addr
	UserAgent string
}

type RotateSessionDTO struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	Role      domain.Role
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
