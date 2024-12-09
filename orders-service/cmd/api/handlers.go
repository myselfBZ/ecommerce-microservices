package main

import (
	"errors"
	"log"
	"net/http"
	"oreders-service/internal/store"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var successMessage = map[string]bool{
	"success": true,
}

var (
	errorBadRequest      = errors.New("bad request")
	errorUserNotFound    = errors.New("user not found")
	errorInternalServer  = errors.New("internal server error")
	errorProductNotFound = errors.New("error product not found")
)

type apiResponse struct {
	Err    string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
	status int
}

type placeOrderRequest struct {
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	Quantity  int    `json"quantity"`
	Address   string `json:"address"`
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

func (a *API) placeOrder(w http.ResponseWriter, r *http.Request) *apiResponse {
	resp := &apiResponse{}
	var orderRequest placeOrderRequest
	if err := readJSON(r, &orderRequest); err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorBadRequest.Error()
		return resp
	}

	_, err := a.userClient.GetByID(orderRequest.UserId)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.status = http.StatusNotFound
			resp.Err = errorUserNotFound.Error()
			return resp
		}
		log.Println("error fetching user: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	prod, err := a.inventoryClient.GetProductById(orderRequest.ProductId)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.status = http.StatusNotFound
			resp.Err = errorProductNotFound.Error()
			return resp
		}
		log.Println("error fetching product: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	if !isEnough(orderRequest.Quantity, int(prod.Quantity)) {
		resp.status = http.StatusNoContent
		resp.Err = "not enough in stock"
		return resp
	}

	inventResp, err := a.inventoryClient.CreateStockTransaction(&pb.StockTransactionRequest{
		ProductId:      int32(orderRequest.ProductId),
		Reason:         "sold",
		Price:          float32(orderRequest.Quantity) * prod.Price,
		QuantityChange: -int32(orderRequest.Quantity),
	})

	if err != nil || !inventResp.Success {
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		log.Println("error creating stock transaction: ", err)
		return resp
	}

	if err := a.store.PlaceOrder(newOrder(&orderRequest)); err != nil {
		log.Println("error creating order record: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	resp.status = http.StatusOK
	resp.Data = successMessage

	return resp

}

func isEnough(orderQnt int, inStock int) bool {
	return inStock >= orderQnt
}

func newOrder(o *placeOrderRequest) *store.Order {
	return &store.Order{
		ProductId:       o.ProductId,
		UserId:          o.UserId,
		ProductQuantity: o.Quantity,
		Address:         o.Address,
	}
}
