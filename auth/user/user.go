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
	Password       string `json:"password,omitempty" gorm:"-"`
	PasswordHashed string `json:"-"`

	OrganizationID uint                 `json:"organization_id" gorm:"index"`
	Organization   *Organization        `json:"organization,omitempty" gorm:"-"`
	Roles          resource.StringSlice `json:"roles"`
	ManagerID      uint                 `json:"manager_id" gorm:"index"`
	Manager        *User                `gorm:"-"`

	Permissions  []*Permission `gorm:"-"`
	Subordinates []uint        `gorm:"-"`
}

func (u *User) Whoami() uint {
	return u.ID
}

func (u *User) Allow(action string, res string, record interface{}) bool {
	for _, p := range u.Permissions {
		if p.Resource == res && p.Action == action {
			return record == nil || p.Record == nil || p.Record.Match(record)
		}
	}
	return false
}

func (u *User) GetFilters(action string, res string) []resource.Filter {
	var filters []resource.Filter
	for _, p := range u.Permissions {
		if p.Action == action && p.Resource == res && p.Record != nil {
			filters = append(filters, p.Record)
		}
	}
	return filters
}

// 1. me 2. org 3. sub (Subordinate)
func (u *User) Convert() {
	permissions := make([]*Permission, 0)
	for _, permission := range u.Permissions {
		isValid := true
		for k, v := range permission.Record {
			switch v {
			case "me":
				permission.Record[k] = u.ID
			case "org":
				if u.OrganizationID == 0 {
					isValid = false
					break
				}
				permission.Record[k] = u.OrganizationID

			case "sub":
				if len(u.Subordinates) > 0 {
					permission.Record[k] = u.Subordinates
				}
			}
		}
		if isValid {
			permissions = append(permissions, permission)
		}
	}
	u.Permissions = permissions
}
