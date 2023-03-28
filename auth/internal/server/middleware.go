package server

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var publicPaths = []string{
	"/auth/register",
	"/auth/login",
}

var resourceExtract = regexp.MustCompile(`\/([^\d^\/]+)\/?(\d*)$`)

type Option struct {
	Secret           string
	PublicPaths      []string
	ResourceExtracts []*regexp.Regexp
}

type authMiddleware struct {
	secret           string
	publicPaths      []string
	resourceExtracts []*regexp.Regexp
}

func (m *authMiddleware) Handle(next middleware.Handler) middleware.Handler {
	return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
		tx, ok := transport.FromServerContext(ctx)
		if !ok {
			return nil, errors.New("Unexcept err when get transport from server context")
		}
		htx := tx.(http.Transporter)
		path := htx.Request().URL.Path
		for _, v := range append(m.publicPaths, publicPaths...) {
			if strings.HasPrefix(path, v) {
				return next(ctx, req)
			}
		}

		var (
			act   string
			res   string
			resID string
		)
		for _, v := range append(m.resourceExtracts, resourceExtract) {
			if find := v.FindStringSubmatch(path); find != nil {
				res, resID = find[1], find[2]
				break
			}
		}
		switch htx.Request().Method {
		case "GET":
			act = "show"
			if resID == "" {
				act = "browse"
			}
		case "PATCH", "POST":
			act = "edit"
			if resID == "" {
				act = "create"
			}
		case "DELETE":
			act = "delete"
		}

		var ac resource.AC
		if !ac.Allow(act, res, nil) {
			return nil, ErrPermissionDenied
		}

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

func MakeAuthMiddleware(opt *Option) middleware.Middleware {
	am := &authMiddleware{secret: opt.Secret, publicPaths: opt.PublicPaths}
	return (middleware.Middleware)(am.Handle)
}

var (
	ErrPermissionDenied = errors.New("permission denied")
)
