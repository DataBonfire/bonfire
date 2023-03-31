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

func (u *User) Allow(action string, _resource string, record interface{}) bool {
	// todo 强制返回 true，测试用
	return true
	for _, p := range u.Permissions {
		if p.Resource == _resource && p.Action == action {
			return true
			//return p.Record == nil || p.Record.Match(record, &resource.UserRelation{
			//	UserId:         u.ID,
			//	OrganizationID: u.OrganizationID,
			//	Subordinates:   u.Subordinates,
			//})
		}
	}
	return false
}

func (u *User) GetFilters(action string, res string) []resource.Filter {

	return []resource.Filter{
		{"created_by": 1},
		{"created_by": []int64{2, 3}, "title": &resource.Constraint{
			Like: "xxxx",
		}},
		{"organization_id": 1},
	}
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
