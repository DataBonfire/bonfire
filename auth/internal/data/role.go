package data

import (
	"context"

	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/go-kratos/kratos/v2/log"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *roleRepo) Find(ctx context.Context, name string) (*user.Role, error) {
	role := user.Role{
		Name: name,
	}
	return &role, r.data.db.First(&role).Error
}
