package data

import (
	"context"
	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/resource"

	"github.com/databonfire/bonfire/auth/user"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (r *userRepo) Find(ctx context.Context, name, email, phone string) (*user.User, error) {
	var data user.User
	return &data, r.data.db.WithContext(ctx).First(&data, "name = ? OR email = ? OR phone = ?", phone).Error
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {

	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}


func (r *userRepo) ACUser(ctx context.Context, id uint) (resource.AC, error) {
	var userInfo user.User
	err := r.data.db.WithContext(ctx).First(&userInfo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	var roles []*user.Role
	err = r.data.db.WithContext(ctx).Where("name in ?", userInfo.Roles).Preload("Permissions").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	permissions := make([]*user.Permission, 0)
	for _, v := range roles {
		permissions = append(permissions, v.Permissions...)
	}
	var users []*user.User
	err = r.data.db.WithContext(ctx).Where("manager_id = ?", userInfo.ID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	subordinates := make([]uint, 0)
	for _, v := range users {
		subordinates = append(subordinates, v.ID)
	}

	userInfo.Subordinates = subordinates
	userInfo.Permissions = permissions

	return &userInfo, nil
}
