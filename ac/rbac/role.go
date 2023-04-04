package rbac

import (
	"strings"

	"github.com/databonfire/bonfire/resource"
)

type Role struct {
	resource.Model
	Name             string        `json:"name" gorm:"type:varchar(60);uniqueIndex"`
	Type             string        `json:"type"`
	IsRegisterPublic bool          `json:"is_register_public"`
	Permissions      []*Permission `json:"-" gorm:"many2many:role_permissions;"`
}

type RoleTemplate struct {
	Name string
	Type string

	ActionsAll            []string
	ActionsMyCreated      []string
	ActionsMySubordinates []string
	ActionsMyCompany      []string

	BrowseAll            []string
	BrowseMyCreated      []string
	BrowseMySubordinates []string
	BrowseMyCompany      []string

	ShowAll            []string
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

	//AdditionalPermissions []*Permission
}

func MakeRoles(tpls ...*RoleTemplate) []*Role {
	var roles []*Role
	for _, t := range tpls {
		var perms []*Permission
		for _, res := range t.ActionsAll {
			perms = append(perms, &Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
			})
		}
		for _, assoicRes := range t.ActionsMyCreated {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, &Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "U"},
			})
		}
		for _, assoicRes := range t.ActionsMySubordinates {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, &Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "S"},
			})
		}
		for _, assoicRes := range t.ActionsMyCompany {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, &Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "C"},
			})
		}
		for _, res := range t.BrowseAll {
			var perm Permission
			perm.Resource = res
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowAll, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.Create, res) {
				perm.Actions = append(perm.Actions, "create")
			}
			if contains(t.EditAll, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteAll, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, &perm)
		}
		for _, assoicRes := range t.BrowseMyCreated {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "U"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowMyCreated, assoicRes) || contains(t.ShowMyCreated, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditMyCreated, assoicRes) || contains(t.EditMyCreated, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteMyCreated, assoicRes) || contains(t.DeleteMyCreated, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, &perm)
		}
		for _, assoicRes := range t.BrowseMySubordinates {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "B"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowMySubordinates, assoicRes) || contains(t.ShowMySubordinates, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditMySubordinates, assoicRes) || contains(t.EditMySubordinates, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteMySubordinates, assoicRes) || contains(t.DeleteMySubordinates, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, &perm)
		}
		for _, assoicRes := range t.BrowseMyCompany {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "C"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowMyCompany, assoicRes) || contains(t.ShowMyCompany, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditMyCompany, assoicRes) || contains(t.EditMyCompany, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteMyCompany, assoicRes) || contains(t.DeleteMyCompany, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, &perm)
		}

		roles = append(roles, &Role{
			Name:        t.Name,
			Type:        t.Type,
			Permissions: perms,
		})
	}
	return roles
}
