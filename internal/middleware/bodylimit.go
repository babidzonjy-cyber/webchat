package middleware

import "net/http"

func MaxbytesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.MaxBytesReader(w, r.Body, 1<<20)
	})
}
