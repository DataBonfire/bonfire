package v1

import (
	"github.com/databonfire/bonfire/ac/rbac"
	http "github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterBlogRBACHTTPServer(s *http.Server, srv BlogHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/posts", rbac.EnhanceContext(_Blog_ListPost0_HTTP_Handler(srv)))
}
