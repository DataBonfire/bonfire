package server

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/databonfire/bonfire/auth/user"

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

		// Access control
		ac, err := getAC(ctx, uid)
		if err != nil {
			return nil, err
		}
		if ac != nil && !ac.Allow(act, res, nil) {
			return nil, ErrPermissionDenied
		}
		ctx = context.WithValue(ctx, "author", ac)
		return next(ctx, req)
	})
}

func getAC(ctx context.Context, uid uint) (resource.AC, error) {
	// 获取 user 基本信息
	userInterface, err := resource.GetRepo("auth.users").Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	userInfo, ok := userInterface.(*user.User)
	if !ok {
		return nil, errors.New("User info error")
	}

	// 根据用户基本信息获取 角色和权限数据
	var roles []*user.Role
	err = resource.GetRepo("auth.roles").DB().Preload("Permissions").
		Where("name in ?", ([]string)(userInfo.Roles)).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	//permissions := make([]*user.Permission, 0)
	//for _, v := range roles {
	//	//permissions = append(permissions, v.Permissions...)
	//}
	//userInfo.Permissions = permissions

	// 根据用户基本信息获取 下属id
	usersInterface, _, err := resource.GetRepo("auth.users").List(ctx, &resource.ListRequest{
		Filter:  resource.Filter{"manager_id": userInfo.ID},
		PerPage: 100,
	})
	subordinates := make([]uint, 0)
	for _, v := range usersInterface {
		_user, ok := v.(*user.User)
		if !ok {
			return nil, errors.New("User info error")
		}
		subordinates = append(subordinates, _user.ID)
	}

	userInfo.Subordinates = subordinates
	return userInfo, nil
}

func MakeAuthMiddleware(opt *Option) middleware.Middleware {
	am := &authMiddleware{secret: opt.Secret, publicPaths: opt.PublicPaths}
	return (middleware.Middleware)(am.Handle)
}

var (
	ErrPermissionDenied = errors.New("permission denied")
)
