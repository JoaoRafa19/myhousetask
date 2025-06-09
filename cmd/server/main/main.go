package main

import (
	"JoaoRafa19/myhousetask/internal/web/handlers"
	m "JoaoRafa19/myhousetask/migrator"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func httpServer() {
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./web/static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", handlers.AdminHandler)

	
	fmt.Println("HTTP server listening on port 8000")

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	db, err := m.Run()
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	defer db.Close()

	l, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	//	pb.RegisterCategoryServiceServer(server, category.NewCategoryServiceServer(db))

	log.Printf("Server listening at %v", l.Addr())

	go httpServer()

	if err := server.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
