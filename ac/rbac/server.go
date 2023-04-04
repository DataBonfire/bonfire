package rbac

import (
	"context"

	"github.com/databonfire/bonfire/ac"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterHTTPServer(srv *http.Server, mws []resource.HTTPHandlerMiddleware) {
	r := srv.Route("/")
	r.GET("/auth/permissions", resource.AssembleHandler(getPermissionsHTTPHandler, mws))
}

func getPermissionsHTTPHandler(ctx http.Context) error {
	reply, err := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
		acer := ctx.Value("acer")
		if acer == nil {
			return nil, ac.ErrUnknownAccessController
		}
		return acer.(*RBAC).Permissions(ctx.Value("author")), nil
	})(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.Result(200, reply)
}
