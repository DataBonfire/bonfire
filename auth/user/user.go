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

	Permissions  []*Permission `gorm:"-"`
	Subordinates []uint        `gorm:"-"`
}

func (u *User) Whoami() uint {
	return u.ID
}

func (u *User) Allow(action string, _resource string, record interface{}) bool {
	// todo 强制返回 true，测试用
	return true
	for _, p := range u.Permissions {
		if p.Resource == _resource && p.Action == action {
			return p.Record == nil || p.Record.Match(record, &resource.UserRelation{
				UserId:         u.ID,
				OrganizationID: u.OrganizationID,
				Subordinates:   u.Subordinates,
			})
		}
	}
	return false
}
