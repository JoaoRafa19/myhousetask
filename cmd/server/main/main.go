package main

import (
	"JoaoRafa19/myhousetask/internal/api"
	"JoaoRafa19/myhousetask/internal/services"
	"JoaoRafa19/myhousetask/store"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

const (
	dbDriver = "mysql"
)

var (
	sm     *scs.SessionManager
	dbUser string
	dbPass string
	dbHost string
	dbPort string
	dbName string
)

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Erro ao carregar arquivo .env: %v\nCarregando variaveis de ambiente do sistema\n", err)
	}
	dbUser = getEnv("DB_USER", "user")
	dbPass = getEnv("DB_PASS", "root")
	dbHost = getEnv("DB_HOST", "localhost")
	dbPort = getEnv("DB_PORT", "3308")
	dbName = getEnv("DB_NAME", "myhousetask")

}

func main() {

	log.Println("Variaveis de ambiente carregadas com sucesso")
	log.Println(dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName)

	// Conecta ao MySQL sem um banco de dados específico para poder criar o nosso.
	initDb, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/")
	if err != nil {
		log.Fatalf("Falha ao conectar no MySQL: %v", err)
	}
	_, err = initDb.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatalf("Falha ao criar o banco de dados: %v", err)
	}
	initDb.Close()

	// Agora, conecta ao banco de dados que sabemos que existe.
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
	if err != nil {
		fmt.Println("Failed start app: ", err)
		os.Exit(1)
	}

	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	sm = scs.New()
	sm = ConfigureSessionManager(db, sm)

	database := store.New(db)
	render := api.NewRenderHandler(database, sm)
	app := api.NewApiHandler(database, sm)

	mux := chi.NewRouter()

	mux.Use(sm.LoadAndSave)
	// Servir arquivos estáticos (CSS, JS)
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))
	FileServer(mux, "/static", filesDir)

	mux.Get("/login", render.LoginPageHandler)

	mux.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Post("/register", app.RegisterUserHandler)
		apiRouter.Post("/login", app.LoginUserHandler)

		apiRouter.Group(func(protectedApiRouter chi.Router) {
			protectedApiRouter.Use(api.AuthRequired(sm))
			protectedApiRouter.Get("/charts/weekly-activity", app.WeeklyActivityHandler)
			protectedApiRouter.Post("/create-family", app.CreateFamilyHandler)
		})
	})

	mux.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(api.AuthRequired(sm))

		protectedRouter.Get("/", render.ShowDashboardPage)
		protectedRouter.Get("/logout", app.LogoutUserHandler)
		protectedRouter.Get("/families", render.ShowFamiliesPage)
		protectedRouter.Get("/tasks", render.ShowTasksPage)

		protectedRouter.Route("/htmx/page", func(htmxRouter chi.Router) {
			htmxRouter.Get("/dashboard", render.RenderDashboardContent)
			htmxRouter.Get("/families", render.RenderFamiliesContent)
			htmxRouter.Get("/families-list", render.RenderFamiliesList)
			//htmxRouter.Get("/users", render.RenderUsersContent)
			htmxRouter.Get("/tasks", render.RenderTasksContent)
		})

		protectedRouter.Route("/htmx", func(htmxRouter chi.Router) {

			htmxRouter.Get("/families-table", render.RenderFamiliesList)
			htmxRouter.Get("/stats-card", render.RenderChart)
			htmxRouter.Get("/dashboard-chart", render.RenderChart)

		})
	})

	fmt.Println("HTTP server listening on port 2345")
	if err := http.ListenAndServe(":2345", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

func ConfigureSessionManager(db *sql.DB, manager *scs.SessionManager) *scs.SessionManager {
	manager.Store = mysqlstore.New(db)
	manager.Lifetime = 24 * time.Hour
	manager.Cookie.Name = services.User_id
	manager.Cookie.HttpOnly = true
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Secure = false

	return manager
}
