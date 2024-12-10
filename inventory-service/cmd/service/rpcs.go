package main

import (
	"context"
	"errors"
	"inventory-service/internal/store"
	"log"
	"time"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	store store.Store
}

func (i *InventoryServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := i.store.GetProductByID(int(r.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	resp := &pb.GetProductResponse{
		ProductId: int32(p.ID),
		Quantity:  int32(p.Quantity),
		Price:     float32(p.Price),
	}

	return resp, nil
}

func (i *InventoryServer) CreateStockTransaction(ctx context.Context, r *pb.StockTransactionRequest) (*pb.StockTransactionResponse, error) {
	newStckTrans := &store.StockTransaction{
		QuantityChange: int(r.QuantityChange),
		Product_id:     int(r.ProductId),
		Reason:         r.Reason,
		CreatedAt:      time.Now(),
	}
	err := i.store.CreateStockTransaction(newStckTrans)
	if err != nil {
		log.Print("error: ", err)
		return nil, errors.New("error database error")
	}
	err = i.store.UpdateProduct(&store.Product{Quantity: int(r.QuantityChange)}, int(r.ProductId))
	if err != nil {
		log.Println("error: ", err)
		return nil, errors.New("error database error")
	}
	return &pb.StockTransactionResponse{
		Success: true,
	}, nil
}
