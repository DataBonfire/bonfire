package resource

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type valueCtx struct {
	http.Context
	key, val any
}

func ContextWithValue(parent http.Context, key, val any) http.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	return &valueCtx{parent, key, val}
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.Context, key)
}

func value(c context.Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
		default:
			return c.Value(key)
		}
	}
}
