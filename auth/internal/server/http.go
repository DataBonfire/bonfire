package server

import (
	"strings"

	stdhttp "net/http"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/service"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, bc *conf.Biz, dc *conf.Data, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			resource.StorageMiddleware,
			MakeAuthMiddleware(&Option{
				Secret: bc.Jwtsecret,
			}),
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
	rdc := &resource.DataConfig{
		Database: &resource.DatabaseConfig{
			Driver: dc.Database.Driver,
			Source: dc.Database.Source,
		},
	}
	resource.RegisterHTTPServer(srv, &resource.Option{
		Resource:   "organizations",
		Model:      &user.Organization{},
		DataConfig: rdc,
		Logger:     logger,
	})
	resource.RegisterHTTPServer(srv, &resource.Option{
		Resource:   "users",
		Model:      &user.User{},
		DataConfig: rdc,
		Logger:     logger,
	})
	resource.RegisterHTTPServer(srv, &resource.Option{
		Resource:   "roles",
		Model:      &user.Role{},
		DataConfig: rdc,
		Logger:     logger,
	})
	resource.RegisterHTTPServer(srv, &resource.Option{
		Resource:   "permissions",
		Model:      &user.Permission{},
		DataConfig: rdc,
		Logger:     logger,
	})
	return srv
}

func readHTTPTransporter(t http.Transporter) (token, path, action, res string) {
	req := t.Request()
	path = req.URL.Path
	token = req.Header.Get("Authorization")
	scheme := "Bearer"
	l := len(scheme)
	if len(token) > l+1 && token[:l] == scheme {
		token = token[l+1:]
	}

	chips := strings.Split(strings.Trim(path, "/"), "/")
	res = chips[0]
	switch req.Method {
	case stdhttp.MethodGet:
		if len(chips) == 1 {
			// GET /posts
			action = resource.ActionBrowse
		} else {
			// GET /posts/1/comments
			action = resource.ActionShow
		}
	case stdhttp.MethodPost, stdhttp.MethodPut, stdhttp.MethodPatch:
		if len(chips) == 1 {
			// [POST|PUT|PATCH] /posts
			action = resource.ActionCreate
		} else {
			// [POST|PUT|PATCH] /posts/1
			// [POST|PUT|PATCH] /posts/1/{action}
			// [POST|PUT|PATCH] /posts/1/comments
			// [POST|PUT|PATCH] /posts/1/comments/1
			action = resource.ActionEdit
		}
	case stdhttp.MethodDelete:
		if len(chips) == 1 {
			// DELETE /posts/1
			action = resource.ActionDelete
		} else {
			// DELETE /posts/1/comments/1
			action = resource.ActionEdit
		}
	}
	return
}
