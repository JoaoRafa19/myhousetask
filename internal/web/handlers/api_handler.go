package handlers

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ApiHandler struct {
	db     *db.Queries
	logger *log.Logger
}

func NewApiHandler(db *db.Queries) *ApiHandler {
	return &ApiHandler{
		db:     db,
		logger: log.Default(),
	}
}

// 1. Defina uma struct para a resposta JSON que o Chart.js espera.
type ChartData struct {
	Labels []string `json:"labels"`
	Data   []int64  `json:"data"`
}

func (h *ApiHandler) WeeklyActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Busca os dados brutos do banco de dados
	activityStats, err := h.db.GetWeeklyTaskCompletionStats(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get weekly activity stats: %v", err)
		http.Error(w, "Failed to get weekly activity stats", http.StatusInternalServerError)
		return
	}

	// 2. Crie uma instância da nossa nova struct de resposta
	chartResponse := ChartData{
		Labels: make([]string, 7),
		Data:   make([]int64, 7),
	}

	for index, stat := range activityStats {
		if stat.CompletionDate.Before(time.Now()) {
			label := stat.CompletionDate.Format("Mon, 02")
			chartResponse.Labels[6-index] = label
			chartResponse.Data[6-index] = stat.CompletedCount
		}
	}

	data, err := json.Marshal(chartResponse)
	if err != nil {
		h.logger.Printf("Failed to marshal weekly activity data: %v", err)
		http.Error(w, "Failed to marshal weekly activity data", http.StatusInternalServerError)
		return
	}

	// Defina o cabeçalho ANTES de escrever o status ou o corpo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ApiHandler) CreateFamilyHandler(w http.ResponseWriter, r *http.Request) {
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
