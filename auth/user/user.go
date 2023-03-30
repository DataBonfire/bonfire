package user

import (
	"github.com/databonfire/bonfire/resource"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	PasswordHashed string `json:"-"`

	OrganizationID uint `json:"organization_id"`
	Organization   *Organization
	Roles          resource.StringSlice `json:"roles"`
	ManagerID      uint                 `json:"manager_id"`
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
