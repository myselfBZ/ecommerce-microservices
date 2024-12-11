package main

import (
	"inventory-service/internal/store"
	"log"
	"net/http"
)

func NewAPI() *API {
	s, err := store.NewPostgre()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	log.Println("connected to the db")
	return &API{
		store: s,
	}
}

type API struct {
	store store.Store
}

// godoc
// @Summary 	mounts all the routes and handlers
func (a *API) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", makeHTTPHandler(a.createProduct, http.MethodPost))
	mux.HandleFunc("PUT /products/{id}", makeHTTPHandler(a.updateProduct, http.MethodPut))
	mux.HandleFunc("GET /products", makeHTTPHandler(a.getProducts, http.MethodGet))
	return mux
}

func (a *API) run(port string) error {
	return http.ListenAndServe(port, a.mount())
}
