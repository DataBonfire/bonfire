package resource

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/middleware"
)

var (
	storage    = make(map[string]Repo)
	storageMtx sync.Mutex
)

func registerRepo(resource string, repo Repo) {
	storageMtx.Lock()
	defer storageMtx.Unlock()
	storage[resource] = repo
}

var StorageMiddleware = func(next middleware.Handler) middleware.Handler {
	return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue("storage", storage)
		return next(ctx, req)
	})
}

//func FromServerContext()
