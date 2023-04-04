package v1

import (
	"github.com/databonfire/bonfire/ac/rbac"
	"github.com/databonfire/bonfire/auth"
	"github.com/databonfire/bonfire/examples/singleton/internal/conf"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/log"
	http "github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterBlogRBACHTTPServer(bc *conf.Biz, s *http.Server, srv BlogHTTPServer, logger log.Logger) {
	r := s.Route("/")
	mws := []resource.HTTPHandlerMiddleware{
		auth.MakeMiddleware(&auth.MiddlewareOption{
			Secret: bc.Jwtsecret,
		}),
		rbac.MakeMiddleware(logger),
	}
	r.GET("/v1/posts", resource.AssembleHandler(_Blog_ListPost0_HTTP_Handler(srv), mws))
}
