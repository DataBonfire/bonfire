package server

import (
	"context"

	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

var makeAuthMiddleware = func(key string) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
			// jwt validate
			var uid uint
			user, err := ctx.Value("storage").(map[string]resource.Repo)["users"].Find(ctx, uid)
			if err != nil {
				return nil, err
			}
			// add ca[*User]
			ctx = context.WithValue(ctx, "author", user)
			return next(ctx, req)
		})
	}
}
