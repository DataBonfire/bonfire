package service

import (
	"context"
	"errors"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
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
	stor := ctx.Value("storage").(map[string]resource.Repo)
	role := &user.Role{Name: req.Role}
	if tx := stor["roles"].DB().First(&role); tx.Error != nil {
		return nil, tx.Error
	}
	if !role.IsRegisterPublic {
		return nil, ErrRegisterIsNotPublic
	}
	u := &user.User{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Password:       req.Password,
		PasswordHashed: hashPassword(req.Password),
		Roles:          []string{req.Role},
	}
	if err := stor["users"].Save(ctx, u); err != nil {
		return nil, err
	}
	return &pb.RegisterReply{}, nil
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	userInfo, tokenStr, err := s.authUsecase.Login(ctx, req.Email, req.Phone, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginReply{
		Token:  tokenStr,
		Name:   userInfo.Name,
		Avatar: userInfo.Avatar,
	}, nil
}

var (
	ErrRegisterIsNotPublic = errors.New("register is not public")
)
