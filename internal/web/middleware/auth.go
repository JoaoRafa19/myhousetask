package middleware

import (
	"JoaoRafa19/myhousetask/internal/core/services"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func AuthRequired(sm *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if sm.GetString(r.Context(), services.User_id) == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Tenta obter o cookie da sessão
			c, err := r.Cookie("myhousetask_session")
			fmt.Println(c)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
