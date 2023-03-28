package user

import (
	"context"

	"github.com/databonfire/bonfire/resource"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string
	Avatar         string
	Email          string
	Phone          string
	Password       string
	PasswordHashed string

	OrganizationID uint
	Organization   *Organization
	Roles          resource.StringSlice
	ManagerID      uint
	Manager        *User

	Permissions []*Permission `gorm:"-"`
}

func (u *User) Whoami() uint {
	return u.ID
}

func (u *User) Allow(action string, resource string, record interface{}) bool {
	for _, p := range u.Permissions {
		if p.Resource == resource && p.Action == action {
			return p.Record == nil || p.Record.Match(record)
		}
	}
	return false
}

type UserRepo interface {
	Save(context.Context, *User) (*User, error)
}
