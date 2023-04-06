package biz

import (
	"context"
	"errors"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
)

type AuthUsecase struct {
	conf     *conf.Biz
	userRepo UserRepo
}

func NewAuthUsecase(c *conf.Biz, userRepo UserRepo) *AuthUsecase {
	return &AuthUsecase{
		c,
		userRepo,
	}
}

func (au *AuthUsecase) Register(ctx context.Context, req *pb.RegisterRequest) error {
	// name, email, phone is duplicate
	_u, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		errMsg := make(map[string]string)
		if _u.Name == req.Name {
			errMsg["name"] = "name is duplicate"
		}
		if _u.Phone == req.Phone {
			errMsg["phone"] = "phone is duplicate"
		}
		if _u.Email == req.Email {
			errMsg["email"] = "email is duplicate"
		}

		return kerrors.BadRequest("account duplicate", "").WithMetadata(errMsg)
	}

	userInfo := &user.User{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		PasswordHashed: utils.HashPassword(req.Password, au.conf.PasswordSalt),
		Roles:          []string{req.Role},
	}
	// save
	err = au.userRepo.Save(ctx, userInfo)
	if err != nil {
		return err
	}

	// notify

	return nil
}

func (au *AuthUsecase) Login(ctx context.Context, req *pb.LoginRequest) (*user.User, string, error) {
	// find user
	userInfo, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil {
		return nil, "", err
	}
	// check password
	passwordHashed := utils.HashPassword(req.Password, au.conf.PasswordSalt)
	if userInfo.PasswordHashed != passwordHashed {
		return nil, "", ErrLoginPassword
	}
	// generate token
	tokenStr, err := utils.GenToken(&utils.UserSession{UserId: userInfo.ID}, au.conf.Jwtsecret)
	if err != nil {
		return nil, "", ErrGenerateToken
	}

	return userInfo, tokenStr, nil
}

var (
	ErrAccountDuplicate    = errors.New("account duplicate")
	ErrRegisterIsNotPublic = errors.New("register is not public")
	ErrLoginPassword       = errors.New("password error")
	ErrGenerateToken       = errors.New("gen token error")
)
