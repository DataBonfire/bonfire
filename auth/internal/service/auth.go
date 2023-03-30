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

	return &pb.RegisterReply{}, s.authUsecase.Register(ctx, req)
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	userInfo, tokenStr, err := s.authUsecase.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token:        tokenStr,
		Name:         userInfo.Name,
		Avatar:       userInfo.Avatar,
		Roles:        userInfo.Roles,
		Organization: &pb.Organization{
			Name: userInfo.Organization.Name,
			Logo: userInfo.Organization.Logo,
		},
	}, nil
}