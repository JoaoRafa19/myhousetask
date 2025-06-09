package main

import (
	m "JoaoRafa19/myhousetask/migrator"
	"log"
	"net"

	"google.golang.org/grpc"
)

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

	if err := server.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
