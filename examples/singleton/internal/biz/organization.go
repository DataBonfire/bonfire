package biz

import (
	"github.com/databonfire/bonfire/auth/user"
	"github.com/databonfire/bonfire/resource"
)

type Organization struct {
	user.Organization
	Industries resource.StringSlice
	Address    string `json:"-"`
}
