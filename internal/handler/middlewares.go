package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

func AuthMiddleware(auth repository.Auth) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				w.WriteHeader(http.StatusUnauthorized)
			}
			session, err := auth.CheckAuthDataBySession(r.Context(), token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", session.UserId)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
