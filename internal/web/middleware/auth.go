package middleware

import (
	"JoaoRafa19/myhousetask/internal/core/services"
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type contextKey string

const userContextKey = contextKey("userID")

func AuthRequired(sm *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := sm.GetString(r.Context(), services.User_id)

			if userID == "" {
				// Utilizador não está logado. Redirecione.
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, userID)

			// 3. Chame o próximo handler com o NOVO contexto.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(userContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}
