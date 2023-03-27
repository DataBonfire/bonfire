package auth

import (
	"context"
	"fmt"
	"unsafe"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Option struct {
	Resources []string
}

var resources = map[string]interface{}{
	"users":         &biz.User{},
	"roles":         &biz.Role{},
	"permissions":   &biz.Permission{},
	"organizations": &biz.Organization{},
}

func RegisterHTTPServer(c *conf.Data, logger log.Logger, srv *http.Server, opt *Option) func() {
	if opt == nil {
		opt = defaultOption()
	}
	authSvc, cleanup, err := wireService(c, logger)
	if err != nil {
		panic(err)
	}
	pb.RegisterAuthHTTPServer(srv, authSvc)
	cleanups := []func(){cleanup}
	for _, res := range opt.Resources {
		if _, ok := resources[res]; !ok {
			panic(fmt.Errorf("Resource %s not found"))
		}
		cleanups = append(cleanups, resource.RegisterHTTPServer((*resource.Config)(unsafe.Pointer(c)), logger, srv, &resource.Option{
			Resource: res,
			Model:    resources[res],
		}))
	}
	return func() {
		for _, f := range cleanups {
			f()
		}
	}
}

func defaultOption() *Option {
	opt := &Option{}
	for res, _ := range resources {
		opt.Resources = append(opt.Resources, res)
	}
	return opt
}

var MakeAuthMiddleware = func(key string) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
			// jwt validate
			var uid uint
			user := ctx.Value("storage").(map[string]resource.Repo)["users"].Find(uid)
			// add ca[*User]
			ctx = context.WithValue("author", user)
			return next(ctx, req)
		})
	}
}
