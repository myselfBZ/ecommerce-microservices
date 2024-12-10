package main

import (
	"log"
	"net"
	"user-service/internal/store"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	store, err := store.NewPostgre()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	pb.RegisterUserServiceServer(server, &UserService{
		store: store,
	})
	ln, err := net.Listen("tcp", "localhost:5423")
	if err != nil {
		log.Fatal("error listening on port 78324: ", err)
	}
	log.Println("service is up")
	if err := server.Serve(ln); err != nil {
		log.Fatal("error on the server: ", err)
	}
}
