package handlers

import (
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	dashboard := pages.DashboardPage()

	dashboard.Render(r.Context(), w)
}
