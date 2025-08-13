package api

import (
	services "JoaoRafa19/myhousetask/internal/services"
	"JoaoRafa19/myhousetask/internal/web/view/components"
	"JoaoRafa19/myhousetask/internal/web/view/pages"
	"JoaoRafa19/myhousetask/store"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type RenderHandler struct {
	db             *store.Queries
	logger         *log.Logger
	sessionManager *scs.SessionManager

	dashboardService *services.DashboardService
	statsCardService *services.StatsCardService
	familyService    *services.FamilyService
	taskService      *services.TaskService
}

func NewRenderHandler(db *store.Queries, sm *scs.SessionManager) *RenderHandler {

	dashboardService := services.NewDashboardService(db)
	statsCardService := services.NewStatsCardService(db)
	familyService := services.NewFamilyService(db)
	taskService := services.NewTaskService(db)

	return &RenderHandler{
		db:               db,
		logger:           log.Default(),
		sessionManager:   sm,
		dashboardService: dashboardService,
		statsCardService: statsCardService,
		familyService:    familyService,
		taskService:      taskService,
	}
}

func (h *RenderHandler) RenderDashboardContent(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	data, _ := h.dashboardService.GetDashboardData(userID)
	pages.DashboardContent(data).Render(r.Context(), w)
	components.Sidebar("dashboard").Render(r.Context(), w)
}

func (h *RenderHandler) ShowDashboardPage(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	h.logger.Println("Utilizando usuario", userID)
	data, _ := h.dashboardService.GetDashboardData(userID)

	pages.DashboardPage(data).Render(r.Context(), w)
	components.Sidebar("dashboard").Render(r.Context(), w)

}

func (h *RenderHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	loginPage := pages.LoginPage()
	err := loginPage.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (h *RenderHandler) RenderFamiliesContent(w http.ResponseWriter, r *http.Request) {
	pages.FamiliesContent().Render(r.Context(), w)
	components.Sidebar("families").Render(r.Context(), w)
}

func (h *RenderHandler) RenderFamiliesList(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

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

func (h *RenderHandler) HTMXStatusCard(w http.ResponseWriter, r *http.Request) {

	userID := GetUserIDFromContext(r.Context())
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

func (h *RenderHandler) ShowFamiliesPage(w http.ResponseWriter, r *http.Request) {
	pages.FamiliesPage().Render(r.Context(), w)
	components.Sidebar("families").Render(r.Context(), w)
}

func (h *RenderHandler) getTasksPageData(r *http.Request) (*services.TasksPageData, error) {
	userID := GetUserIDFromContext(r.Context())

	families, err := h.familyService.GetFamiliesByUserID(r.Context(), userID)
	if err != nil || len(families) == 0 {
		return nil, fmt.Errorf("could not fetch families: %w", err)
	}

	// Por enquanto, seleciona a primeira família por padrão.
	var selectedFamily store.Family = families[0]

	tasks, err := h.taskService.GetTasksByFamily(r.Context(), selectedFamily.ID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch tasks: %w", err)
	}

	taskStatus, err := h.taskService.GetTaskStatus(r.Context())
	if err != nil {
		return nil, fmt.Errorf("could not fetch task status: %w", err)
	}

	return &services.TasksPageData{
		Families:       families,
		Tasks:          tasks,
		TaskStatus:     taskStatus,
		SelectedFamily: selectedFamily,
	}, nil
}

func (h *RenderHandler) ShowTasksPage(w http.ResponseWriter, r *http.Request) {
	data, err := h.getTasksPageData(r)
	if err != nil {
		h.logger.Printf("Error getting tasks page data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	pages.TasksPage(data).Render(r.Context(), w)
}

func (h *RenderHandler) RenderTasksContent(w http.ResponseWriter, r *http.Request) {
	data, err := h.getTasksPageData(r)
	if err != nil {
		h.logger.Printf("Error getting tasks page data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	pages.TasksContent(data).Render(r.Context(), w)
	components.Sidebar("tasks").Render(r.Context(), w)
}
