package handlers

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/core/services"
	"JoaoRafa19/myhousetask/internal/web/view/components"
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type RenderHandler struct {
	db             *db.Queries
	logger         *log.Logger
	sessionManager *scs.SessionManager

	dashboardService *services.DashboardService
	statsCardService *services.StatsCardService
}

func NewRenderHandler(db *db.Queries, sm *scs.SessionManager) *RenderHandler {

	dashboardService := services.NewDashboardService(db)
	if dashboardService == nil {
		log.Fatal("Failed to create DashboardService")
	}
	statsCardService := services.NewStatsCardService(db)
	if statsCardService == nil {
		log.Fatal("Failed to create StatsCardService")
	}

	return &RenderHandler{
		db:               db,
		logger:           log.Default(),
		dashboardService: dashboardService,
		sessionManager:   sm,
		statsCardService: statsCardService,
	}
}

func (h *RenderHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {

	userID := h.sessionManager.GetString(r.Context(), "userID")
	h.logger.Println("Utilizando usuario", userID)
	data, err := h.dashboardService.GetDashboardData(userID)
	if err != nil {
		log.Printf("Error getting dashboard data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(data)

	dashboard := pages.DashboardPage(data)

	dashboard.Render(r.Context(), w)
}

func (h *RenderHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Render the login page
	loginPage := pages.LoginPage()
	loginPage.Render(r.Context(), w)
}

func (h *RenderHandler) FamiliesTableHTMXHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling HTMX request for families table")
	families, err := h.db.ListRecentFamilies(r.Context())
	if err != nil {
		log.Printf("Erro ao buscar famílias: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	table := pages.FamiliesTableComponent(families)
	table.Render(r.Context(), w)
}

func (h *RenderHandler) HTMXStatusCard(w http.ResponseWriter, r *http.Request) {

	userID := h.sessionManager.GetString(r.Context(), "userID")
	if userID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data, err := h.statsCardService.GetStatsCardData(r.Context(), userID)

	if err != nil {
		log.Printf("Error getting stats card data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	statsCard := components.StatsCards(data)
	statsCard.Render(r.Context(), w)

}
