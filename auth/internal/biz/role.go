package biz

import (
	"context"

	"github.com/databonfire/bonfire/auth/user"
)

type RoleRepo interface {
	Find(context.Context, string) (*user.Role, error)
}
