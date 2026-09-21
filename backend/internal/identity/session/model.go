package session

import (
	"net/netip"
	"time"
	"uuid"
)

type Session struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	TokenHash  string    `db:"token_hash"`
	DeviceInfo string    `db:"device_info"`
	DeviceHash string    `db:"device_hash"`

	IPAddress     netip.Addr `db:"ip_address"`
	LastIPAddress netip.Addr `db:"last_ip_address"`

	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
