package session

import (
	"net/netip"
	"time"
	"uuid"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:user_sessions,alias:s"`

	ID            uuid.UUID  `bun:"id,pk,type:text"`
	UserID        uuid.UUID  `bun:"user_id,type:text"`
	TokenHash     string     `bun:"token_hash"                                    json:"-"`
	DeviceInfo    string     `bun:"device_info"`
	IPAddress     netip.Addr `bun:"ip_address"`
	LastIPAddress netip.Addr `bun:"last_ip_address"`
	ExpiresAt     time.Time  `bun:"expires_at"`
	CreatedAt     time.Time  `bun:"created_at,nullzero,default:current_timestamp"`
}
