package biz

import (
	"context"
	"errors"
	"os/user"

	"github.com/databonfire/bonfire/auth/internal/conf"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	conf     *conf.Biz
	roleRepo RoleRepo
	userRepo UserRepo
}

func NewAuthUsecase(c *conf.Biz, roleRepo RoleRepo, userRepo UserRepo) *AuthUsecase {
	return &AuthUsecase{
		c,
		roleRepo,
		userRepo,
	}
}

func (au *AuthUsecase) Register(ctx context.Context, u *user.User) error {
	// role is public register?
	r, err := au.roleRepo.Find(ctx, u.Role)
	if err != nil {
		return err
	}
	if !r.IsRegisterPublic {
		return Err
	}

	// name, email, phone is duplicate
	_, err = au.userRepo.Find(ctx, u.Name, u.Email, u.Phone)
	if err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		return ErrAccountDuplicate
	}
	// encrypt password
	// save
	// notice
}

var (
	ErrAccountDuplicate    = errors.New("account duplicate")
	ErrRegisterIsNotPublic = errors.New("register is not public")
)
