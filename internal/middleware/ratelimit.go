package middleware

import "net/http"

type RateLimiter struct {
	sem chan struct{}
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case rl.sem <- struct{}{}:
			defer func() { <-rl.sem }()
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "too many requests", http.StatusTooManyRequests)
		}
	})
}
