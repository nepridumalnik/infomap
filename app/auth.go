package app

import (
	"net/http"
)

type middleware struct {
	storage *storage
}

func (m *middleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
