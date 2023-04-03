package server

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
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
	Secret             string
	PublicPaths        []string
	ResourceExtracts   []*regexp.Regexp
	SubordinatesFinder func(ctx context.Context, u *user.User) ([]uint, error)
}

type authMiddleware struct {
	secret             string
	publicPaths        []string
	resourceExtracts   []*regexp.Regexp
	subordinatesFinder func(ctx context.Context, u *user.User) ([]uint, error)
}

func (m *authMiddleware) Handle(next middleware.Handler) middleware.Handler {
	return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
		tx, ok := transport.FromServerContext(ctx)
		if !ok {
			return nil, errors.New("Unexcept err when get transport from server context")
		}
		htx := tx.(http.Transporter)
		token, path := readHTTPTransporter(htx)

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
		db := resource.GetRepo("auth.users").(resource.Repo).DB()
		var u user.User
		if err = db.First(&u, userSession.UserId).Error; err != nil {
			return nil, err
		}
		if m.subordinatesFinder != nil {
			u.Subordinates, err = m.subordinatesFinder(ctx, &u)
		} else {
			err = db.Model(&user.User{}).Select("id").Where("manager_id", u.ID).Find(&u.Subordinates).Error
		}
		if err != nil {
			return nil, err
		}
		return next(context.WithValue(ctx, "author", &u), req)
	})
}

func MakeAuthMiddleware(opt *Option) middleware.Middleware {
	am := &authMiddleware{secret: opt.Secret, publicPaths: opt.PublicPaths, subordinatesFinder: opt.SubordinatesFinder}
	return (middleware.Middleware)(am.Handle)
}
