package main

import (
	"log"
	"net/http"

	"user-service/internal/store"
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
	mux.HandleFunc("/register", makeHTTPHandler(a.registerUser, http.MethodPost))
	mux.HandleFunc("/login", makeHTTPHandler(a.login, http.MethodPost))
	mux.HandleFunc("DELETE /users/{id}", makeHTTPHandler(a.deleteAccount, http.MethodDelete))
	mux.HandleFunc("PUT /users/{id}", makeHTTPHandler(a.updateAccount, http.MethodPut))
	return mux
}

func (a *API) run(port string) error {
	return http.ListenAndServe(port, a.mount())
}
