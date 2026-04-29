package handler

import (
	"context"
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

func AuthMiddleware(auth repository.Auth) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sessionID string
			sessionID = r.Header.Get("Authorization")

			if sessionID == "" {
				cookie, err := r.Cookie("session_id")
				if err == nil && cookie != nil {
					sessionID = cookie.Value
				}
			}

			if sessionID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			session, err := auth.CheckAuthDataBySession(r.Context(), sessionID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", session.UserId)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
