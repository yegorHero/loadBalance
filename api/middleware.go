package api

import (
	"context"
	"loadBalance/internal/rateLimited"
	"net"
	"net/http"
	"sync"
	"time"
)

func RateLimitedMiddleware(ctx context.Context, rate time.Duration, capacity int) func(http.Handler) http.Handler {
	limiters := sync.Map{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				responseWithJSON(w, http.StatusInternalServerError, "Unable to parse IP")
				return
			}

			limiterIface, _ := limiters.LoadOrStore(ip, rateLimited.NewTokenBucket(ctx, rate, capacity))
			limiter := limiterIface.(*rateLimited.TokenBucket)

			if !limiter.Allow() {
				responseWithJSON(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
