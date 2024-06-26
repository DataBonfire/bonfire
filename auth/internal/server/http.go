package server

import (
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
			resource.Validator(),
		),
		http.Filter(resource.MakeCors()),
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
		Resource: user.OrganizationResourceName,
		Model:    &user.Organization{},
		HTTPHandlerMiddlewares: []resource.HTTPHandlerMiddleware{
			MakeAuthMiddleware(&Option{
				Secret: bc.Jwtsecret,
			}),
		},
		DataConfig: rdc,
		Logger:     logger,
	})
	resource.RegisterHTTPServer(srv, &resource.Option{
		Resource: "users",
		Model:    &user.User{},
		HTTPHandlerMiddlewares: []resource.HTTPHandlerMiddleware{
			MakeAuthMiddleware(&Option{
				Secret: bc.Jwtsecret,
			}),
		},
		DataConfig: rdc,
		Logger:     logger,
	})
	//resource.RegisterHTTPServer(srv, &resource.Option{
	//	Resource:   "roles",
	//	Model:      &user.Role{},
	//	DataConfig: rdc,
	//	Logger:     logger,
	//})
	//resource.RegisterHTTPServer(srv, &resource.Option{
	//	Resource:   "permissions",
	//	Model:      &user.Permission{},
	//	DataConfig: rdc,
	//	Logger:     logger,
	//})
	return srv
}

func readHTTPTransporter(t http.Transporter) (token, path string, info *requestInfo) {
	req := t.Request()
	info = &requestInfo{}
	info.Method = req.Method
	info.IP = getAccessIP(req)
	info.UserAgent = req.Header.Get("User-Agent")
	path = req.URL.Path
	token = req.Header.Get("Authorization")
	scheme := "Bearer"
	l := len(scheme)
	if len(token) > l+1 && token[:l] == scheme {
		token = token[l+1:]
	}
	return
}

type requestInfo struct {
	Method        string
	IP            string
	UserAgent     string
	RequestedWith string
	Referer       string
}

func getAccessIP(c *http.Request) string {
	ks := []string{"X-Real-IP", "X-Forwarded-For"}
	for _, v := range ks {
		accessIP := c.Header.Get(v)
		if len(accessIP) != 0 {
			return accessIP
		}
	}
	return ""
}
