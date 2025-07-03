package main

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/web/handlers"
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

	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))

	FileServer(mux, "/static", filesDir)

	mux.Get("/", render.DashboardHandler)

	mux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/charts/weekly-activity", api.WeeklyActivityHandler)
			r.Post("/create-family", api.CreateFamilyHandler)
		})
	})
	
	mux.Route("/htmx", func(r chi.Router) {
		r.Get("/families-table", render.FamiliesTableHTMXHandler)
		r.Get("/stats-card", render.HTMXStatusCard)
	})


	fmt.Println("HTTP server listening on port 3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// FileServer convenientemente serve arquivos estáticos de um diretório.
// Esta é uma função helper recomendada pela própria documentação do Chi.
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
