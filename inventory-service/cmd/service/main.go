package main

import (
	"log"
	"net"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterInventoryServiceServer(server, &InventoryServer{})
	ln, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("error listening: ", err)
	}

	log.Println("inventory service is running on port 50052")
	if err := server.Serve(ln); err != nil {
		log.Fatal("error: ", err)
	}
}
