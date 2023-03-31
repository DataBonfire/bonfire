package service

import (
	"context"
	"encoding/json"
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
	reply := &pb.LoginReply{
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

func (s *AuthService) GetPermissions(ctx context.Context, req *pb.GetPermissionsRequest) (*pb.GetPermissionsReply, error) {
	permissions, err := s.authUsecase.GetPermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	getPermissionsReply := &pb.GetPermissionsReply{Permissions: make([]*pb.GetPermissionsReply_Permission, 0)}
	for _, permission := range permissions {
		if permission.Record == nil {
			continue
		}
		tempBytes, err := json.Marshal(permission)
		if err != nil {
			return nil, err
		}
		var p pb.GetPermissionsReply_Permission
		if err = json.Unmarshal(tempBytes, &p); err != nil {
			return nil, err
		}

		getPermissionsReply.Permissions = append(getPermissionsReply.Permissions, &p)
	}

	return getPermissionsReply, nil
}
