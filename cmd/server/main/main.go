package main

import (
	db "JoaoRafa19/myhousetask/db/db/gen"
	"JoaoRafa19/myhousetask/internal/web/handlers"
	m "JoaoRafa19/myhousetask/migrator"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	mydb, err := m.Run()

	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	defer mydb.Close()

	database := db.New(mydb)

	handler := handlers.NewHandler(database)

	router.HandleFunc("POST /create-family", handler.CreateFamilyHandler)

	fs := http.FileServer(http.Dir("./web/static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", handlers.AdminHandler)

	fmt.Println("HTTP server listening on port 8000")

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
