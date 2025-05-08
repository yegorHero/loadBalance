package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func RateLimitedMiddleware(rate time.Duration, capacity int) func(http.Handler) http.Handler {
	limiters := sync.Map{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Unable to parse IP", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
