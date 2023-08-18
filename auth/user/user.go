package user

import (
	"github.com/databonfire/bonfire/resource"
)

type User struct {
	resource.Model
	Type           string `json:"type"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	Phone          string `json:"phone"`
	PhoneVerified  bool   `json:"phone_verified"`
	Password       string `json:"password,omitempty" gorm:"-"`
	PasswordHashed string `json:"-"`

	OrganizationID uint                 `json:"organization_id" gorm:"index"`
	Organization   *Organization        `json:"organization,omitempty" gorm:"-"`
	Roles          resource.StringSlice `json:"roles"`
	ManagerID      uint                 `json:"manager_id" gorm:"index"`
	Manager        *User                `gorm:"-"`

	Subordinates []uint `gorm:"-"`
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetGroups() []uint {
	if u.OrganizationID > 0 {
		return []uint{u.OrganizationID}
	}
	return nil
}

func (u *User) GetSubordinates() []uint {
	return u.Subordinates
}

func (u *User) GetRoles() []string {
	return u.Roles
}

func (u *User) GetRoleType() string {
	return u.Type
}
