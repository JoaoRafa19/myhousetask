package handlers

import (
	db "JoaoRafa19/myhousetask/db/db/gen"
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"database/sql"
	"net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	dashboard := pages.DashboardPage()

	dashboard.Render(r.Context(), w)
}

func (h *Handler) CreateFamilyHandler(w http.ResponseWriter, r *http.Request) {
	family := r.FormValue("familyName")
	description := r.FormValue("familyDescription")

	if family == "" {
		http.Error(w, "Family name is required", http.StatusBadRequest)
		return
	}

	err := h.db.CreateFamily(r.Context(), db.CreateFamilyParams{
		Name:        family,
		Description: sql.NullString{String: description, Valid: true},
	})

	if err != nil {
		h.logger.Printf("Failed to create family: %v", err)
		http.Error(w, "Failed to create family", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Family: %s, Description: %s", family, description)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
