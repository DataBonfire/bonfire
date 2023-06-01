package server

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"strings"

	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type authHandler struct {
	opt  *Option
	next stdhttp.Handler
}

func (h *authHandler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	token, path := readHTTPRequest(r)

	// Public Endpoints
	for _, v := range append(h.opt.PublicPaths, publicPaths...) {
		if strings.HasPrefix(path, v) {
			h.next.ServeHTTP(w, r)
			return
		}
	}

	// Auth
	userSession, err := utils.ParseToken(token, h.opt.Secret)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := resource.GetRepo("auth.users").(resource.Repo).DB()
	var u user.User
	if err = db.First(&u, userSession.UserId).Error; err != nil {
		fmt.Println(err)
		return
	}
	if h.opt.SubordinatesFinder != nil {
		u.Subordinates, err = h.opt.SubordinatesFinder(r.Context(), &u)
	} else {
		err = db.Model(&user.User{}).Select("id").Where("manager_id", u.ID).Find(&u.Subordinates).Error
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.WithValue(r.Context(), "author", &u)
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func MakeAuthFilter(opt *Option) http.FilterFunc {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return &authHandler{
			opt,
			next,
		}
	}
}
