package user

import "context"

type HookFunc func(context.Context, *User) error

const (
	ON_REGISTER_EMAIL_VERIFY = "on_register_email_verify"
	ON_REGISTER_PHONE_VERIFY = "on_register_phone_verify"
	ON_REGISTER_SUCCESS      = "on_register_success"
	ON_REGISTER_VERIFIED     = "on_register_verified"

	ON_LOGIN_EMAIL_VERIFY = "on_login_email_verify"
	ON_LOGIN_PHONE_VERIFY = "on_login_phone_verify"

	ON_FORGET_PASSWORD = "on_forget_password"
)
