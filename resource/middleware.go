package resource

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

var StorageMiddleware = func(next middleware.Handler) middleware.Handler {
	return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, "storage", storage)
		return next(ctx, req)
	})
}
