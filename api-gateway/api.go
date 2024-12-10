package main

import (
	"io"
	"net/http"
	"os"
)

var (
	ORDERS_SERVICE    = os.Getenv("orders")
	USERS_SERVICE     = os.Getenv("users")
	INVENTORY_SERVICE = os.Getenv("inventory")
)

type API struct {
	middleware []Middleware
}

func NewAPI() *API {
	return &API{}
}

func forwardHeaders(resp *http.Response, w http.ResponseWriter) {
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
}

func (a *API) handleRequests(w http.ResponseWriter, r *http.Request) {

	services := map[string]string{
		"/orders":   ORDERS_SERVICE,
		"/products": INVENTORY_SERVICE,
		"/users":    USERS_SERVICE,
	}

	service, ok := services[r.URL.Path]

	if !ok {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	req, err := http.NewRequest(r.Method, service+r.URL.Path, r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	forwardHeaders(resp, w)
	io.Copy(w, resp.Body)
	w.WriteHeader(resp.StatusCode)
}

func (a *API) use(m Middleware) {
	a.middleware = append(a.middleware, m)
}

func (a *API) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleRequests)
	//middleware happens here
	a.use(JWTValidate)
	handler := middlewareFunc(mux)
	return handler
}

func (a *API) run(addr string) error {
	return http.ListenAndServe(addr, a.mount())
}
