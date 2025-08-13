package api

import (
	"JoaoRafa19/myhousetask/internal/services"
	"JoaoRafa19/myhousetask/store"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ApiHandler struct {
	db     *store.Queries
	logger *log.Logger
	sm     *scs.SessionManager
}

func NewApiHandler(db *store.Queries, sm *scs.SessionManager) *ApiHandler {
	return &ApiHandler{
		db:     db,
		logger: log.Default(),
		sm:     sm,
	}
}

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
	if _, err := w.Write(data); err != nil {
		h.logger.Printf("Failed to write weekly activity data: %v", err)
		http.Error(w, "Failed to write weekly activity data", http.StatusInternalServerError)
		return
	}
}

func (h *ApiHandler) CreateFamilyHandler(w http.ResponseWriter, r *http.Request) {
	family := r.FormValue("familyName")
	description := r.FormValue("familyDescription")

	if family == "" {
		http.Error(w, "Family name is required", http.StatusBadRequest)
		return
	}

	// Get the currently logged-in user ID from session
	userID := h.sm.GetString(r.Context(), services.User_id)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Create the family with the current user as owner
	result, err := h.db.CreateFamily(r.Context(), store.CreateFamilyParams{
		Name:        family,
		Description: sql.NullString{String: description, Valid: true},
		OwnerID:     sql.NullString{String: userID, Valid: true},
	})

	if err != nil {
		h.logger.Printf("Failed to create family: %v", err)
		http.Error(w, "Failed to create family", http.StatusInternalServerError)
		return
	}

	// Get the inserted family ID
	familyID, err := result.LastInsertId()
	if err != nil {
		h.logger.Printf("Failed to get family ID: %v", err)
		http.Error(w, "Failed to get family ID", http.StatusInternalServerError)
		return
	}

	// Add the creator as a family member with OWNER role
	err = h.db.CreateFamilyMember(r.Context(), store.CreateFamilyMemberParams{
		ID:       uuid.New().String(),
		FamilyID: sql.NullInt32{Int32: int32(familyID), Valid: true},
		UserID:   sql.NullString{String: userID, Valid: true},
		Role:     sql.NullString{String: "OWNER", Valid: true},
	})

	if err != nil {
		h.logger.Printf("Failed to add user as family owner: %v", err)
		http.Error(w, "Failed to add user as family owner", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Family created: %s (ID: %d), Description: %s, Owner: %s", family, familyID, description, userID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ApiHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		http.Error(w, "As senhas não coincidem", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao fazer hash da senha: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	err = h.db.CreateUser(r.Context(), store.CreateUserParams{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		http.Error(w, "Não foi possível criar o usuário", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *ApiHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse dos dados do formulário
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 2. Buscar o usuário pelo email
	fmt.Println("get user:", email)
	user, err := h.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Usuário não encontrado - mensagem genérica por segurança
			log.Println("Tentativa de login falhou (usuário não encontrado):", email)
			http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
			return
		}
		// Outro erro de banco de dados
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// 3. Comparar a senha fornecida com o hash armazenado
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		// Senha incorreta - bcrypt.ErrMismatchedHashAndPassword
		log.Println("Tentativa de login falhou (senha incorreta):", email)
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	err = h.sm.RenewToken(r.Context())
	if err != nil {
		http.Error(w, "Erro ao renovar token", http.StatusInternalServerError)
		return
	}

	h.sm.Put(r.Context(), services.User_id, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ApiHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {

	err := h.sm.Destroy(r.Context())
	if err != nil {
		http.Error(w, "Erro ao logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
