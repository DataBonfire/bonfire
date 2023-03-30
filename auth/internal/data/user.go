package data

import (
	"context"

	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
	//return &data, r.data.db.WithContext(ctx).Preload("Organization").First(&data, "name = ? OR email = ? OR phone = ?", name, email, phone).Error
	if tx := r.data.db.WithContext(ctx).First(&data, "name = ? OR email = ? OR phone = ?", name, email, phone); tx.Error != nil {
		return nil, tx.Error
	}
	if data.OrganizationID == 0 {
		return &data, nil
	}
	data.Organization = &user.Organization{Model: gorm.Model{ID: data.OrganizationID}}
	return &data, r.data.db.First(data.Organization).Error
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	return r.data.db.WithContext(ctx).Save(u).Error
}
