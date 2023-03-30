package biz

import (
	"context"
	"github.com/databonfire/bonfire/auth/user"
)

type UserRepo interface {
	// Name, Email, Phone
	Find(ctx context.Context, name, email, phone string) (*user.User, error)
	Save(context.Context, *user.User) error
}
