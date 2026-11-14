package middleware

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

type ipCtxKey struct{}

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := ExtractIP(r)
		if err != nil {
			writeAuthError(w, http.StatusBadRequest, apperr.CodeBadRequest, "invalid IP address")
			return
		}
		ctx := context.WithValue(r.Context(), ipCtxKey{}, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIP(ctx context.Context) (netip.Addr, bool) {
	val, ok := ctx.Value(ipCtxKey{}).(netip.Addr)
	return val, ok
}

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
