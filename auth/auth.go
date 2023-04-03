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
	Resources              map[string]interface{}
	HTTPHandlerMiddlewares []resource.HTTPHandlerMiddleware
	DataConfig             *resource.DataConfig
	JWTSecret              string
	PasswordSalt           string
	PublicRegisterRoles    []string
	Logger                 log.Logger
}

func defaultResources() map[string]interface{} {
	return map[string]interface{}{
		"users": &user.User{},
		//"roles":         &user.Role{},
		//"permissions":   &user.Permission{},
		user.OrganizationResourceName: &user.Organization{},
	}
}

func RegisterHTTPServer(srv *http.Server, opt *Option) func() {
	authSvc, cleanup, err := wireService(&conf.Biz{
		Jwtsecret:           opt.JWTSecret,
		PasswordSalt:        opt.PasswordSalt,
		PublicRegisterRoles: opt.PublicRegisterRoles,
	}, &conf.Data{
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
			Resource:               k,
			Model:                  v,
			DataConfig:             opt.DataConfig,
			Logger:                 opt.Logger,
			HTTPHandlerMiddlewares: opt.HTTPHandlerMiddlewares,
		}))
	}
	for k, v := range defaultResources() {
		//if _, ok := opt.Resources[k]; !ok {
		//	continue
		//}
		if cleanup := resource.RegisterHTTPServer(srv, &resource.Option{
			AuthPackage: true,
			Resource:    k,
			Model:       v,
			DataConfig:  opt.DataConfig,
			Logger:      opt.Logger,
		}); cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	return func() {
		for _, f := range cleanups {
			f()
		}
	}
}

type MiddlewareOption server.Option

func MakeMiddleware(opt *MiddlewareOption) middleware.Middleware {
	return server.MakeAuthMiddleware((*server.Option)(opt))
}
