package api

import (
	"JoaoRafa19/myhousetask/internal/services"
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type contextKey string

const userContextKey = contextKey("userID")

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(services.User_id).(string)
	if !ok {
		return ""
	}
	return userID
}

func AuthRequired(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sm.GetString(r.Context(), services.User_id)
			if userID == "" {
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), services.User_id, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
