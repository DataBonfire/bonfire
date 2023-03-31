package rbac

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name             string
	IsRegisterPublic bool
	Permissions      []*Permission `gorm:"many2many:role_permissions;"`
}
