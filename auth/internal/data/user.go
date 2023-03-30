package data

import (
	"context"

	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Find(ctx context.Context, name, email, phone string) (*user.User, error) {
	var data user.User
	return &data, r.data.db.WithContext(ctx).Preload("Organization").First(&data, "name = ? OR email = ? OR phone = ?", name, email, phone).Error
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	return r.data.db.WithContext(ctx).Save(u).Error
}

