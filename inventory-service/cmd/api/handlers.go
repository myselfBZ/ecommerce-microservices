package main

import (
	"database/sql"
	"errors"
	"inventory-service/internal/store"
	"log"
	"net/http"
	"strconv"
)

// Errors
var (
	errorInternalServer = errors.New("internal server error")
	errorBadRequest     = errors.New("bad request")
	errorNotFound       = errors.New("product not found")
)

var successMessage = map[string]bool{
	"success": true,
}

type productCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type apiResponse struct {
	Err    string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
	status int
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiResponse

func makeHTTPHandler(f apiHandler, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiResp := f(w, r)
		if apiResp.Err != "" {
			writeJSON(w, apiResp, apiResp.status)
			return
		}
		writeJSON(w, apiResp, apiResp.status)
	}
}

func (a *API) createProduct(w http.ResponseWriter, r *http.Request) *apiResponse {
	resp := &apiResponse{}
	var p productCreateRequest
	if err := readJSON(r, &p); err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorBadRequest.Error()
		return resp
	}

	productStore := newStoreProduct(&p)

	if err := a.store.CreateProduct(productStore); err != nil {
		log.Println("error creating product: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	resp.status = http.StatusOK
	resp.Data = successMessage
	return resp

}

func (a *API) updateProduct(w http.ResponseWriter, r *http.Request) *apiResponse {
	var resp = &apiResponse{}
	id := r.PathValue("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorBadRequest.Error()
		return resp
	}
	var p productCreateRequest
	if err := readJSON(r, &p); err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorBadRequest.Error()
		return resp
	}

	storeProd := newStoreProduct(&p)

	if err := a.store.UpdateProduct(storeProd, validId); err != nil {
		if err == sql.ErrNoRows {
			resp.status = http.StatusNotFound
			resp.Err = errorNotFound.Error()
			return resp
		}
		log.Println("error updating a product: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	resp.status = http.StatusOK
	resp.Data = successMessage
	return resp
}

func newStoreProduct(p *productCreateRequest) *store.Product {
	return &store.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
	}
}
