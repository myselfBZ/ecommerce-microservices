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

func newApiResponse(err string, data any, status int) *apiResponse {
	return &apiResponse{
		status: status,
		Data:   data,
		Err:    err,
	}
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

	var p productCreateRequest
	if err := readJSON(r, &p); err != nil {
		return newApiResponse(errorBadRequest.Error(), nil, http.StatusBadRequest)
	}

	productStore := newStoreProduct(&p)

	if err := a.store.CreateProduct(productStore); err != nil {
		log.Println("error creating product: ", err)
		return newApiResponse(errorInternalServer.Error(), nil, http.StatusInternalServerError)
	}

	return newApiResponse("", successMessage, http.StatusOK)

}

func (a *API) updateProduct(w http.ResponseWriter, r *http.Request) *apiResponse {
	id := r.PathValue("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		return newApiResponse(errorBadRequest.Error(), nil, http.StatusBadRequest)
	}
	var p productCreateRequest
	if err := readJSON(r, &p); err != nil {
		return newApiResponse(errorBadRequest.Error(), nil, http.StatusBadRequest)
	}

	storeProd := newStoreProduct(&p)

	if err := a.store.UpdateProduct(storeProd, validId); err != nil {
		if err == sql.ErrNoRows {
			return newApiResponse(errorNotFound.Error(), nil, http.StatusNotFound)
		}
		log.Println("error updating a product: ", err)
		return newApiResponse(errorInternalServer.Error(), nil, http.StatusInternalServerError)
	}

	return newApiResponse("", successMessage, http.StatusOK)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) *apiResponse {
	products, err := a.store.GetProducts()
	if err != nil {
		return newApiResponse(errorInternalServer.Error(), nil, http.StatusInternalServerError)
	}
	return newApiResponse("", products, http.StatusOK)
}

func newStoreProduct(p *productCreateRequest) *store.Product {
	return &store.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
	}
}
