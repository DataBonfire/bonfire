package rbac

import (
	"context"
	stdhttp "net/http"

	"github.com/databonfire/bonfire/ac"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type handler struct {
	rbac   ac.AccessController
	logger *log.Helper
	next   stdhttp.Handler
}

func (h *handler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var a Accessor
	if v := r.Context().Value("author"); v != nil {
		if accessor, ok := v.(Accessor); ok {
			a = accessor
		}
	}

	act, res := ac.ReadHTTPRequest(r)

	if res != "auth" && !h.rbac.Allow(a, act, res, nil) {
		h.logger.Error(ac.ErrPermissionDenied)
		return
	}
	ctx := context.WithValue(r.Context(), "acer", h.rbac)
	//ctx := &Context{http.Wrapper{}}
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func MakeFilter(logger log.Logger) http.FilterFunc {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return &handler{
			newAC(nil, logger),
			log.NewHelper(logger),
			next,
		}
	}
}
