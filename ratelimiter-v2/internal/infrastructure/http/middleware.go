package http

import "net/http"

func RatelimitterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here

		// IP request/s limit

		next.ServeHTTP(w, r)
	})
}
