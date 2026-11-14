package middleware

import (
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

func ExtractIP(r *http.Request) (netip.Addr, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	remoteIP, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, apperr.New(apperr.CodeBadRequest, "invalid IP address")
	}
	remoteIP = remoteIP.Unmap()

	if remoteIP.IsLoopback() || remoteIP.IsPrivate() {
		if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
			if parsed, err := netip.ParseAddr(strings.TrimSpace(cfIP)); err == nil {
				return parsed.Unmap().WithZone(""), nil
			}
		}

		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			if parsed, err := netip.ParseAddr(strings.TrimSpace(realIP)); err == nil {
				return parsed.Unmap().WithZone(""), nil
			}
		}

		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			firstIP, _, _ := strings.Cut(xff, ",")
			if parsed, err := netip.ParseAddr(strings.TrimSpace(firstIP)); err == nil {
				return parsed.Unmap().WithZone(""), nil
			}
		}
	}

	return remoteIP.WithZone(""), nil
}
