//go:build wireinject
// +build wireinject

package auth

import (
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/data"
	"github.com/databonfire/bonfire/auth/internal/service"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireService(*conf.Biz, *conf.Data, log.Logger, map[string]user.HookFunc) (*service.AuthService, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.NewAuthService))
}
