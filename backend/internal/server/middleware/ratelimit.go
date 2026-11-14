package middleware

import (
	"net"
	"net/http"
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
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			now := time.Now()
			mu.Lock()
			v, exists := visitors[ip]
			if !exists || now.After(v.expiresAt) {
				visitors[ip] = &visitor{count: 1, expiresAt: now.Add(window)}
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
