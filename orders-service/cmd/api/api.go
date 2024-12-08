package main

import (
	"log"
	"net/http"
	grpc "oreders-service/internal/gRPC"
	"oreders-service/internal/store"
)

func NewAPI() *API {
	s, err := store.NewPostgre()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	log.Println("connected to the db")
	usrClient := grpc.NewUserClient("localhost:50051")
	inventoryClient := grpc.NewInventoryClient("localhost:784324")
	return &API{
		store:           s,
		userClient:      usrClient,
		inventoryClient: inventoryClient,
	}
}

type API struct {
	store           store.Store
	userClient      *grpc.UserClient
	inventoryClient *grpc.InventoryClient
}

// godoc
// @Summary 	mounts all the routes and handlers
func (a *API) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", makeHTTPHandler(a.placeOrder, http.MethodPost))
	return mux
}

func (a *API) run(port string) error {
	return http.ListenAndServe(port, a.mount())
}
