package server

import (
	"context"
	"errors"
	v1 "github.com/databonfire/bonfire/auth/api/v1"
	"regexp"
	"strings"
	"time"

	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var publicPaths = []string{
	"/auth/register",
	"/auth/login",
	"/auth/forget_password",
	"/auth/reset_password",
}

type Option struct {
	Secret             string
	PublicPaths        []string
	ResourceExtracts   []*regexp.Regexp
	SubordinatesFinder func(ctx context.Context, u *user.User) ([]uint, error)

	AccessLogIgnores []string
}

type authMiddleware struct {
	secret             string
	publicPaths        []string
	resourceExtracts   []*regexp.Regexp
	subordinatesFinder func(ctx context.Context, u *user.User) ([]uint, error)

	accessLogIgnores []string
}

func (m *authMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return (http.HandlerFunc)(func(ctx http.Context) error {
		tx, ok := transport.FromServerContext(ctx)
		if !ok {
			return errors.New("Unexcept err when get transport from server context")
		}
		htx := tx.(http.Transporter)
		token, path, _requestInfo := readHTTPTransporter(htx)

		// Public Endpoints
		for _, v := range append(m.publicPaths, publicPaths...) {
			if strings.HasPrefix(path, v) {
				return next(ctx)
			}
		}

		// Auth
		userSession, err := utils.ParseToken(token, m.secret)
		if err != nil {
			return ErrNeedLogin
		}
		db := resource.GetRepo("auth.users").(resource.Repo).DB()
		var u user.User
		if err = db.First(&u, userSession.UserId).Error; err != nil {
			return err
		}
		if u.ExpiredAt <= time.Now().Unix() {
			return ErrNeedLogin
		}
		if u.SecureUpdatedAt > 0 {
			if userSession.TokenIssuedAt <= u.SecureUpdatedAt {
				return ErrNeedLogin
			}
		}

		if m.subordinatesFinder != nil {
			u.Subordinates, err = m.subordinatesFinder(ctx, &u)
		} else {
			err = db.Model(&user.User{}).Select("id").Where("manager_id", u.ID).Find(&u.Subordinates).Error
		}
		if err != nil {
			return err
		}
		_userLog := &user.UserLog{
			UserId:     u.ID,
			MethodPath: _requestInfo.Method,
			AccessAt:   time.Now().Unix(),
			UA:         _requestInfo.UserAgent,
			IP:         _requestInfo.IP,
			RequestUri: path,
		}
		if u.ID > 0 {
			ignore := false
			for _, v := range []string{"sts_tokens", "global_unread_messages", "influencer_comparisons", "influencer_collections"} {
				if strings.Contains(path, v) {
					ignore = true
					break
				}
			}
			if !ignore {
				if errCreate := db.Create(_userLog).Error; errCreate != nil {
					println(errCreate.Error())
					return nil
				}
			}
		}

		return next(resource.ContextWithValue(ctx, "author", &u))
	})
}

func MakeAuthMiddleware(opt *Option) resource.HTTPHandlerMiddleware {
	am := &authMiddleware{secret: opt.Secret, publicPaths: opt.PublicPaths, subordinatesFinder: opt.SubordinatesFinder, accessLogIgnores: opt.AccessLogIgnores}
	return (resource.HTTPHandlerMiddleware)(am.Handle)
}

var (
	//ErrNeedLogin = kerrors.Unauthorized(v1.ErrorReason_NEED_LOGIN.String(), "need login")
	ErrNeedLogin = v1.ErrorNeedLogin("need login")
)
