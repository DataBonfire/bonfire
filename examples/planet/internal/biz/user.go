package biz

import "github.com/databonfire/bonfire/auth/user"

type User struct {
	user.User
	YOB        string
	Avatar     string
	Position   string
	Biography  string
	IsVerified bool
}
