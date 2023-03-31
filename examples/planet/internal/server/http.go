package server

import (
	"github.com/databonfire/bonfire/auth"
	v1 "github.com/databonfire/bonfire/examples/planet/api/planet/v1"
	"github.com/databonfire/bonfire/examples/planet/internal/biz"
	"github.com/databonfire/bonfire/examples/planet/internal/conf"
	"github.com/databonfire/bonfire/examples/planet/internal/service"
	"github.com/databonfire/bonfire/resource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, bc *conf.Biz, dc *conf.Data, blog *service.InfluencerService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			auth.MakeAuthMiddleware(bc.Jwtsecret, nil),
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
	v1.RegisterInfluencerHTTPServer(srv, blog)
	{
		rdc := &resource.DataConfig{
			Database: &resource.DatabaseConfig{
				Driver: dc.Database.Driver,
				Source: dc.Database.Source,
			},
		}
		auth.RegisterHTTPServer(srv, &auth.Option{
			Resources: map[string]interface{}{
				"users":         &biz.User{},
				"organizations": &biz.Organization{},
				"influencers":   &biz.Influencer{},
				//"posts.comments": &biz.Comment{},
				//"posts.comments.replies": &biz.Reply{},
			},
			DataConfig:   rdc,
			JWTSecret:    bc.Jwtsecret,
			PasswordSalt: bc.PasswordSalt,
			Logger:       logger,
		})
		//resource.RegisterHTTPServer(srv, &resource.Option{
		//	Resource:   "posts",
		//	Model:      &biz.Post{},
		//	DataConfig: rdc,
		//	Logger:     logger,
		//})
	}
	return srv
}
