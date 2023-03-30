package auth

import (
	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/server"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Option struct {
	Resources  map[string]interface{}
	DataConfig *resource.DataConfig
	Logger     log.Logger
}

var resources = map[string]interface{}{
	"users":         &user.User{},
	"roles":         &user.Role{},
	"permissions":   &user.Permission{},
	"organizations": &user.Organization{},
}

func RegisterHTTPServer(srv *http.Server, opt *Option) func() {
	authSvc, cleanup, err := wireService(&conf.Biz{}, &conf.Data{
		Database: &conf.Data_Database{
			Driver: opt.DataConfig.Database.Driver,
			Source: opt.DataConfig.Database.Source,
		},
	}, opt.Logger)
	if err != nil {
		panic(err)
	}
	pb.RegisterAuthHTTPServer(srv, authSvc)

	cleanups := []func(){cleanup}
	for k, v := range opt.Resources {
		cleanups = append(cleanups, resource.RegisterHTTPServer(srv, &resource.Option{
			Resource:   k,
			Model:      v,
			DataConfig: opt.DataConfig,
			Logger:     opt.Logger,
		}))
	}
	for k, v := range resources {
		if _, ok := opt.Resources[k]; !ok {
			continue
		}
		cleanups = append(cleanups, resource.RegisterHTTPServer(srv, &resource.Option{
			Resource:   k,
			Model:      v,
			DataConfig: opt.DataConfig,
			Logger:     opt.Logger,
		}))
	}

	return func() {
		for _, f := range cleanups {
			f()
		}
	}
}

func MakeAuthMiddleware(secret string, publicPaths []string) middleware.Middleware {
	return server.MakeAuthMiddleware(&server.Option{
		Secret:      secret,
		PublicPaths: publicPaths,
	})
}
