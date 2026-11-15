package middleware

import (
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

type visitor struct {
	count     int
	expiresAt time.Time
}

func RateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	var (
		mu          sync.Mutex
		visitors    = make(map[netip.Addr]*visitor)
		lastCleanup = time.Now()
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, ok := GetIP(r.Context())

			if !ok || !ip.IsValid() {
				if extracted, err := ExtractIP(r); err == nil && extracted.IsValid() {
					ip = extracted
				} else {
					writeAuthError(
						w,
						http.StatusBadRequest,
						apperr.CodeBadRequest,
						"unable to determine client IP",
					)
					return
				}
			}

			now := time.Now()

			mu.Lock()

			if now.Sub(lastCleanup) >= window {
				for addr, v := range visitors {
					if now.After(v.expiresAt) {
						delete(visitors, addr)
					}
				}

				lastCleanup = now
			}

			v, exists := visitors[ip]
			if !exists || now.After(v.expiresAt) {
				visitors[ip] = &visitor{
					count:     1,
					expiresAt: now.Add(window),
				}

				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if v.count >= maxRequests {
				mu.Unlock()
				writeAuthError(
					w,
					http.StatusTooManyRequests,
					apperr.CodeTooManyRequests,
					"Too many requests",
				)
				return
			}

			v.count++
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
