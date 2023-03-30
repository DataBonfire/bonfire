// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package auth

import (
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/data"
	"github.com/databonfire/bonfire/auth/internal/service"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

func wireService(confData *conf.Data, logger log.Logger) (*service.AuthService, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	authUsecase := biz.NewAuthUsecase(userRepo)
	authService := service.NewAuthService(authUsecase)
	return authService, func() {
		cleanup()
	}, nil
}