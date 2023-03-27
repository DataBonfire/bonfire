package server

import (
	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/service"
	"github.com/databonfire/bonfire/resource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, dc *conf.Data, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			//auth.AC(dc),
			//auth.RemoteAC(dc),
			resource.StorageMiddleware,
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	pb.RegisterAuthHTTPServer(srv, auth)
	rc := &resource.Config{
		&resource.Database{
			Driver: dc.Database.Driver,
			Source: dc.Database.Source,
		},
		nil,
	}
	auth.RegisterHTTPServer(rc, logger, srv)
	resource.RegisterHTTPServer(rc, logger, srv, &resource.Option{
		Resource: "organizations",
		Model:    &biz.Organization{},
	})
	resource.RegisterHTTPServer(rc, logger, srv, &resource.Option{
		Resource: "users",
		Model:    &biz.User{},
	})
	resource.RegisterHTTPServer(rc, logger, srv, &resource.Option{
		Resource: "roles",
		Model:    &biz.Role{},
	})
	resource.RegisterHTTPServer(rc, logger, srv, &resource.Option{
		Resource: "permissions",
		Model:    &biz.Permission{},
	})
	return srv
}
