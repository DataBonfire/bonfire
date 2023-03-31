package rbac

import "github.com/databonfire/bonfire/resource"

type Role struct {
	resource.Model
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	IsRegisterPublic bool          `json:"is_register_public"`
	Permissions      []*Permission `json:"-" gorm:"many2many:role_permissions;"`
}
