package handlers

import (
	db "JoaoRafa19/myhousetask/db/db/gen"
	"JoaoRafa19/myhousetask/internal/core/services"
	"log"
)

type Handler struct {
	db     *db.Queries
	logger *log.Logger

	dashboardService *services.DashboardService
}

func NewHandler(db *db.Queries) *Handler {

	dashboardService := services.NewDashboardService(db)

	return &Handler{db: db, logger: log.Default(), dashboardService: dashboardService}
}
