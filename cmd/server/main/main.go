package main

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/web/handlers"
	"JoaoRafa19/myhousetask/internal/web/middleware" // Certifique-se de importar o seu middleware
	m "JoaoRafa19/myhousetask/migrator"
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
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

	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))
	FileServer(mux, "/static", filesDir)

	mux.Use(sessionManager.LoadAndSave)
	mux.Get("/login", render.LoginPageHandler)

	mux.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Post("/register", api.RegisterUserHandler)
		apiRouter.Post("/login", api.LoginUserHandler)

		apiRouter.Group(func(protectedApiRouter chi.Router) {
			protectedApiRouter.Use(middleware.AuthRequired)
			protectedApiRouter.Get("/charts/weekly-activity", api.WeeklyActivityHandler)
			protectedApiRouter.Post("/create-family", api.CreateFamilyHandler)
		})
	})

	mux.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.AuthRequired)

		protectedRouter.Get("/", render.DashboardHandler)

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
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
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
	manager.Cookie.Name = "myhousetask_session"
	manager.Cookie.HttpOnly = true
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Secure = false

	return manager
}
