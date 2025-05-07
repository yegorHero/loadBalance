package api

import (
	"net/http"
	"sync"
	"time"
)

func RateLimitedMiddleware(rate time.Duration, capacity int) func(http.Handler) http.Handler {
	limiters := sync.Map{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	}
}
