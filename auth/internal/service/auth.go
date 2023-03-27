package service

import (
	"context"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServer

	authUsecase *biz.AuthUsecase
}

func NewAuthService(au *biz.AuthUsecase) *AuthService {
	return &AuthService{
		authUsecase: au,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
