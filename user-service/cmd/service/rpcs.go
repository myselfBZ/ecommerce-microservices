package main

import (
	"context"
	"errors"
	"user-service/internal/store"

	pb "github.com/myselfBZ/common-grpc/pkg"
)

var userNotFound = errors.New("user not found")

type UserService struct {
	pb.UnimplementedUserServiceServer
	store store.Store
}

func (s *UserService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.store.GetById(int(r.UserId))
	if err != nil {
		return nil, userNotFound
	}
	return &pb.GetUserResponse{
		Name:   u.Name,
		UserId: int32(u.ID),
		Email:  u.Email,
	}, nil
}
