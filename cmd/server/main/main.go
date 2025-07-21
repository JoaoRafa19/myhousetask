package main

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/core/services"
	"JoaoRafa19/myhousetask/internal/web/handlers"
	"JoaoRafa19/myhousetask/internal/web/middleware" // Certifique-se de importar o seu middleware
	m "JoaoRafa19/myhousetask/migrator"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

var sessionManager *scs.SessionManager

func main() {
	myDb, err := m.Run()
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	defer func() {
		if err := myDb.Close(); err != nil {
			log.Fatalf("could not close db: %v", err)
		}
	}()

	sessionManager = scs.New()
	sessionManager = ConfigureSessionManager(myDb, sessionManager)

	database := db.New(myDb)
	render := handlers.NewRenderHandler(database, sessionManager)
	api := handlers.NewApiHandler(database, sessionManager)

	mux := chi.NewRouter()

	mux.Use(sessionManager.LoadAndSave)
	// Servir arquivos estáticos (CSS, JS)
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))
	FileServer(mux, "/static", filesDir)

	mux.Get("/login", render.LoginPageHandler)

	mux.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Post("/register", api.RegisterUserHandler)
		apiRouter.Post("/login", api.LoginUserHandler)

		apiRouter.Group(func(protectedApiRouter chi.Router) {
			protectedApiRouter.Use(middleware.AuthRequired(sessionManager))
			protectedApiRouter.Get("/charts/weekly-activity", api.WeeklyActivityHandler)
			protectedApiRouter.Post("/create-family", api.CreateFamilyHandler)
		})
	})

	mux.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.AuthRequired(sessionManager))

		protectedRouter.Get("/", render.DashboardHandler)
		protectedRouter.Get("/logout", api.LogoutUserHandler)
		protectedRouter.Get("/families", render.ShowFamiliesPage)

		protectedRouter.Route("/htmx", func(htmxRouter chi.Router) {

			htmxRouter.Get("/families-table", render.FamiliesTableHTMXHandler)
			htmxRouter.Get("/stats-card", render.HTMXStatusCard)
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
