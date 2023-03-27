package biz

import "github.com/databonfire/bonfire/auth/user"

type AuthUsecase struct {
	userRepo user.UserRepo
}

func NewAuthUsecase(userRepo user.UserRepo) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}
