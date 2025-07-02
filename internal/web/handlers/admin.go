package handlers

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"database/sql"
	"log"
	"net/http"
)

func (h *Handler) AdminHandler(w http.ResponseWriter, r *http.Request) {

	data, err := h.dashboardService.GetDashboardData()
	if err != nil {
		log.Printf("Error getting dashboard data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(data)

	dashboard := pages.DashboardPage(&data)

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
