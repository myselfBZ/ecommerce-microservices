package grpc

import (
	"context"
	"log"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient struct {
	inventoryClient pb.InventoryServiceClient
}

func NewInventoryClient(addr string) *InventoryClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting inventory service: ", err)
	}
	return &InventoryClient{
		inventoryClient: pb.NewInventoryServiceClient(conn),
	}
}

func (i *InventoryClient) GetProductById(id int) (*pb.GetProductResponse, error) {
	resp, err := i.inventoryClient.GetProduct(context.TODO(), &pb.GetProductRequest{ProductId: int32(id)})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *InventoryClient) CreateStockTransaction(r *pb.StockTransactionRequest) (*pb.StockTransactionResponse, error) {
	resp, err := i.inventoryClient.CreateStockTransaction(context.TODO(), r)
	return resp, err
}
