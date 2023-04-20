package user

import "context"

type HookFunc func(context.Context, *User) error

const (
	ON_REGISTER_SUCCESS = "on_register_success"
)
