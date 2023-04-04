package v1

import (
	"github.com/databonfire/bonfire/ac/rbac"
	"github.com/databonfire/bonfire/auth"
	"github.com/databonfire/bonfire/resource"
	http "github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterBlogRBACHTTPServer(s *http.Server, srv BlogHTTPServer) {
	r := s.Route("/")
	mws := []resource.HTTPHandlerMiddleware{
		auth.MakeMiddleware(&auth.MiddlewareOption{
			Secret: bc.Jwtsecret,
		}),
		rbac.MakeMiddleware(logger),
		rbac.EnhanceContext,
	}
	r.GET("/v1/posts", resource.AssembleHandler(_Blog_ListPost0_HTTP_Handler(srv), mws))
}
