package rbac

import (
	"context"
	"errors"

	"github.com/databonfire/bonfire/ac"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func MakeACMiddleware(logger log.Logger) middleware.Middleware {
	rbac := newAC(nil, logger)
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var a Accessor
			if v := ctx.Value("author"); v != nil {
				if accessor, ok := v.(Accessor); ok {
					a = accessor
				}
			}

			tx, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New("Unexcept err when get transport from server context")
			}
			htx := tx.(http.Transporter)
			act, res := ac.ReadHTTPTransporter(htx)

			if res != "auth" && !rbac.Allow(a, act, res, nil) {
				return nil, ac.ErrPermissionDenied
			}

			return next(context.WithValue(ctx, "acer", rbac), req)
		}
	}
}
