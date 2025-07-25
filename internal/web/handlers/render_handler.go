package handlers

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/core/services"
	"JoaoRafa19/myhousetask/internal/web/middleware"
	"JoaoRafa19/myhousetask/internal/web/view/components"
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"fmt"
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
	familyService    *services.FamilyServices
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

	familyService := services.NewFamilyServices(db)
	if familyService == nil {
		log.Fatal("Failed to create FamilyService")
	}

	return &RenderHandler{
		db:               db,
		logger:           log.Default(),
		dashboardService: dashboardService,
		sessionManager:   sm,
		statsCardService: statsCardService,
		familyService:    familyService,
	}
}

func (h *RenderHandler) RenderDashboardContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	data, _ := h.dashboardService.GetDashboardData(userID)
	pages.DashboardContent(data).Render(r.Context(), w)
	components.Sidebar("dashboard").Render(r.Context(), w)
}

func (h *RenderHandler) ShowDashboardPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	h.logger.Println("Utilizando usuario", userID)
	data, _ := h.dashboardService.GetDashboardData(userID)

	pages.DashboardPage(data).Render(r.Context(), w)
	components.Sidebar("dashboard").Render(r.Context(), w)

}

func (h *RenderHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserIDFromContext(r.Context())
	h.logger.Println("Utilizando usuario", userID)
	data, err := h.dashboardService.GetDashboardData(userID)
	if err != nil {
		log.Printf("Error getting dashboard data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(data)

	dashboard := pages.DashboardPage(data)

	err = dashboard.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (h *RenderHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Render the login page
	loginPage := pages.LoginPage()
	err := loginPage.Render(r.Context(), w)
	if err != nil {
		return
	}
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

	userID := middleware.GetUserIDFromContext(r.Context())
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

func (h *RenderHandler) RenderFamiliesContent(w http.ResponseWriter, r *http.Request) {
	pages.FamiliesContent().Render(r.Context(), w)
	components.Sidebar("families").Render(r.Context(), w)
}

func (h *RenderHandler) RenderUsersContent(writer http.ResponseWriter, request *http.Request) {

}

func (h *RenderHandler) RenderFamiliesList(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	if userID == "" {
		http.Error(w, "Utilizador não autorizado", http.StatusUnauthorized)
		return
	}

	families, err := h.familyService.GetFamiliesByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Printf("Erro ao buscar famílias para o utilizador %s: %v", userID, err)
		http.Error(w, "Não foi possível carregar as famílias", http.StatusInternalServerError)
		return
	}

	pages.FamiliesList(families).Render(r.Context(), w)
}

func (h *RenderHandler) RenderChart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RENDERIZANDO O GRÁFICO")
	components.ActivityChart().Render(r.Context(), w)
}
