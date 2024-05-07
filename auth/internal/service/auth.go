package service

import (
	"context"
	"fmt"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/internal/conf"
)

type AuthService struct {
	pb.UnimplementedAuthServer

	authUsecase         *biz.AuthUsecase
	publicRegisterRoles []string
}

func NewAuthService(c *conf.Biz, au *biz.AuthUsecase) *AuthService {
	return &AuthService{
		authUsecase:         au,
		publicRegisterRoles: c.PublicRegisterRoles,
	}
}

func (s *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.CommonReply, error) {

	return &pb.CommonReply{Id: "1"}, s.authUsecase.ResetPassword(ctx, req)
}

func (s *AuthService) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) (*pb.CommonReply, error) {
	return &pb.CommonReply{Id: "1"}, s.authUsecase.ForgetPassword(ctx, req)
}

func (s *AuthService) ResendRegister(ctx context.Context, req *pb.ResendRegisterRequest) (*pb.CommonReply, error) {
	return &pb.CommonReply{Id: "1"}, s.authUsecase.ResendRegister(ctx, req)
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	var registerValid bool
	for _, v := range s.publicRegisterRoles {
		if v == req.Role {
			registerValid = true
			break
		}
	}
	if !registerValid {
		return nil, fmt.Errorf("%s is not public register.", req.Role)
	}
	return &pb.RegisterReply{}, s.authUsecase.Register(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	userInfo, tokenStr, err := s.authUsecase.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pb.LoginReply{
		Id:     uint32(userInfo.ID),
		Token:  tokenStr,
		Name:   userInfo.Name,
		Avatar: userInfo.Avatar,
		Roles:  userInfo.Roles,
	}
	if userInfo.Organization != nil {
		reply.Organization = &pb.Organization{
			Name: userInfo.Organization.Name,
			Logo: userInfo.Organization.Logo,
		}
	}
	return reply, nil
}
