package rbac

import (
	"errors"

	"github.com/databonfire/bonfire/ac"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func MakeMiddleware(logger log.Logger) resource.HTTPHandlerMiddleware {
	rbac := newAC(nil, logger)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) error {
			var a Accessor
			if v := ctx.Value("author"); v != nil {
				if accessor, ok := v.(Accessor); ok {
					a = accessor
				}
			}

			tx, ok := transport.FromServerContext(ctx)
			if !ok {
				return errors.New("Unexcept err when get transport from server context")
			}
			htx := tx.(http.Transporter)
			act, res := ac.ReadHTTPTransporter(htx)

			if res != "auth" && !rbac.Allow(a, act, res, nil) {
				return ac.ErrPermissionDenied
			}
			return next(resource.ContextWithValue(ctx, "acer", rbac))
		}
	}
}
