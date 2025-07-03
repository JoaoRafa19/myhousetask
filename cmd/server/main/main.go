package main

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/web/handlers"
	"JoaoRafa19/myhousetask/internal/web/middleware" // Certifique-se de importar o seu middleware
	m "JoaoRafa19/myhousetask/migrator"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

)

func main() {
	mydb, err := m.Run()
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	defer mydb.Close()

	database := db.New(mydb)
	render := handlers.NewRenderHandler(database)
	api := handlers.NewApiHandler(database)

	mux := chi.NewRouter()

	// Servir arquivos estáticos (CSS, JS)
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))
	FileServer(mux, "/static", filesDir)


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
		// (Adicione outras páginas aqui, ex: /families, /users)

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
