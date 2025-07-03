package middleware

import (
	"net/http"
)

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Tenta obter o cookie da sessão
		_, err := r.Cookie("myhousetask_session")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
