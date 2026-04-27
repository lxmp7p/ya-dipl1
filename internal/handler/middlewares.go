package handler

import (
	"net/http"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//
		})
	}
}
