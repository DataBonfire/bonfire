package rbac

import "github.com/databonfire/bonfire/resource"

type Role struct {
	resource.Model
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	IsRegisterPublic bool          `json:"is_register_public"`
	Permissions      []*Permission `json:"-" gorm:"many2many:role_permissions;"`
}

type RoleTemplate struct {
	Name string
	Type string

	Browse               []string
	BrowseMyCreated      []string
	BrowseMySubordinates []string
	BrowseMyCompany      []string

	Show               []string
	ShowMyCreated      []string
	ShowMySubordinates []string
	ShowMyCompany      []string

	Create []string

	EditAll            []string
	EditMyCreated      []string
	EditMySubordinates []string
	EditMyCompany      []string

	DeleteAll            []string
	DeleteMyCreated      []string
	DeleteMySubordinates []string
	DeleteMyCompany      []string

	AdditionalPermissions []*Permission
}

func MakeRoles(tpls []*RoleTemplate) []*Role {
	var roles []*Role
	for _, v := range tpls {
		//var all, my, sub, com []string
		roles = append(roles, &Role{
			Name: v.Name,
			Type: v.Type,
		})
	}
	return roles
}
