package grpc

import (
	"commerce/cmd/user/usecase"
	"commerce/proto/userpb"
	"context"
)

type GRPCServer struct {
	userpb.UnimplementedUserServiceServer
	UserUseCase usecase.UserUseCase
}

func (s *GRPCServer) GetUserByUserId(ctx context.Context, request *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {
	userInfo, err := s.UserUseCase.GetUserById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &userpb.GetUserInfoResponse{
		UserId: userInfo.ID,
		Name:   userInfo.Name,
		Email:  userInfo.Email,
		Role:   userInfo.Role,
	}, nil
}
