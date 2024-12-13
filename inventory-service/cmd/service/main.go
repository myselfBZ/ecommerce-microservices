package main

import (
	"inventory-service/internal/store"
	"log"
	"net"
	"os"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	store, err := store.NewPostgre()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	pb.RegisterInventoryServiceServer(server, &InventoryServer{
		store: store,
	})
	ln, err := net.Listen("tcp", os.Getenv("inventory-service"))
	if err != nil {
		log.Fatal("error listening: ", err)
	}

	log.Println("inventory service is running on port 50052")
	if err := server.Serve(ln); err != nil {
		log.Fatal("error: ", err)
	}
}
