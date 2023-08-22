package data

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/databonfire/bonfire/auth/internal/biz"
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
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
	if len(name) == 0 && len(email) == 0 && len(phone) == 0 {
		return nil, errors.New("all params is empty")
	}

	chains := r.data.db.WithContext(ctx).Where("")
	if len(name) != 0 {
		chains.Or("name = ?", name)
	}
	if len(email) != 0 {
		chains.Or("email = ?", email)
	}
	if len(phone) != 0 {
		chains.Or("phone = ?", phone)
	}

	if err := chains.Limit(1).Find(&data).Error; err != nil {
		return nil, err
	}
	if data.ID <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if data.OrganizationID == 0 {
		return &data, nil
	}
	data.Organization = &user.Organization{Model: resource.Model{ID: data.OrganizationID}}
	return &data, r.data.db.First(data.Organization).Error
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	return r.data.db.WithContext(ctx).Save(u).Error
}
