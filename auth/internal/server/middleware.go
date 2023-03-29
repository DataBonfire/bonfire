package server

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/databonfire/bonfire/auth/internal/utils"

	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var publicPaths = []string{
	"/auth/register",
	"/auth/login",
}

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
		token, path, act, res := readHTTPTransporter(htx)

		// Public Endpoints
		for _, v := range append(m.publicPaths, publicPaths...) {
			if strings.HasPrefix(path, v) {
				return next(ctx, req)
			}
		}

		// Auth
		userSession, err := utils.ParseToken(token, m.secret)
		if err != nil {
			return nil, err
		}
		uid := userSession.UserId

		// todo AC
		// user -> roles
		// roles -> permissions
		// match

		// Access control
		var ac resource.AC
		user, err := ctx.Value("storage").(map[string]resource.Repo)["users"].Find(ctx, uid)
		if err != nil {
			return nil, err
		}
		_ = user
		//ac := user.AC()
		if ac != nil && !ac.Allow(act, res, nil) {
			return nil, ErrPermissionDenied
		}
		ctx = context.WithValue(ctx, "author", ac)
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
