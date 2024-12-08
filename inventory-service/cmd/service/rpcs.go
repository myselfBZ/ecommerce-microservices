package main

import (
	"context"
	"inventory-service/internal/store"

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

func (i *InventoryServer) CreateStocktransaction(ctx context.Context, r *pb.StockTransactionRequest) (*pb.StockTransactionResponse, error) {
	return nil, nil
}
