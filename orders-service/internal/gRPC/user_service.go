package grpc

import (
	"context"
	"log"

	pb "github.com/myselfBZ/common-grpc/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(addr string) *UserClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting to users service")
	}
	return &UserClient{
		client: pb.NewUserServiceClient(conn),
	}
}

func (u *UserClient) GetByID(id int) (*pb.GetUserResponse, error) {
	resp, err := u.client.GetUser(context.TODO(), &pb.GetUserRequest{UserId: int32(id)})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
