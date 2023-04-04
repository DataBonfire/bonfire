package rbac

import (
	"bytes"
	"fmt"
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

	ActionsAll          []string
	ActionsUID          []string
	ActionsSubordinates []string
	ActionsCompany      []string

	BrowseAll          []string
	BrowseUID          []string
	BrowseSubordinates []string
	BrowseCompany      []string

	ShowAll          []string
	ShowUID          []string
	ShowSubordinates []string
	ShowCompany      []string

	Create []string

	EditAll          []string
	EditUID          []string
	EditSubordinates []string
	EditCompany      []string

	DeleteAll          []string
	DeleteUID          []string
	DeleteSubordinates []string
	DeleteCompany      []string

	//AdditionalPermissions []*Permission
}

func MakeRoles(tpls []*RoleTemplate) []*Role {
	var (
		roles      []*Role
		perms      = map[string]uint{}
		bindPermID = func(perm *Permission) *Permission {
			var buf bytes.Buffer
			buf.WriteString(perm.Resource)
			buf.WriteString(strings.Join(perm.Actions, ","))
			for k, v := range perm.Record {
				buf.WriteString(fmt.Sprintf("%s=%v", k, v))
			}
			key := buf.String()
			id := perms[key]
			if id == 0 {
				id = uint(len(perms) + 1)
				perms[key] = id
			}
			perm.ID = id
			return perm
		}
	)

	for _, t := range tpls {
		var perms []*Permission
		for _, res := range t.ActionsAll {
			perms = append(perms, bindPermID(&Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
			}))
		}
		for _, assoicRes := range t.ActionsUID {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, bindPermID(&Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "U"},
			}))
		}
		for _, assoicRes := range t.ActionsSubordinates {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, bindPermID(&Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "S"},
			}))
		}
		for _, assoicRes := range t.ActionsCompany {
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perms = append(perms, bindPermID(&Permission{
				Resource: res,
				Actions:  []string{"browse", "show", "create", "edit", "delete"},
				Record:   map[string]interface{}{assoic: "C"},
			}))
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

			perms = append(perms, bindPermID(&perm))
		}
		for _, assoicRes := range t.BrowseUID {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "U"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowUID, assoicRes) || contains(t.ShowUID, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditUID, assoicRes) || contains(t.EditUID, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteUID, assoicRes) || contains(t.DeleteUID, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, bindPermID(&perm))
		}
		for _, assoicRes := range t.BrowseSubordinates {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "B"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowSubordinates, assoicRes) || contains(t.ShowSubordinates, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditSubordinates, assoicRes) || contains(t.EditSubordinates, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteSubordinates, assoicRes) || contains(t.DeleteSubordinates, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, bindPermID(&perm))
		}
		for _, assoicRes := range t.BrowseCompany {
			var perm Permission
			v := strings.Split(assoicRes, ".")
			if len(v) != 2 {
				panic("unexcepted resource or associated field")
			}
			res, assoic := v[0], v[1]
			perm.Resource = res
			perm.Record = map[string]interface{}{assoic: "C"}
			perm.Actions = append(perm.Actions, "browse")
			if contains(t.ShowCompany, assoicRes) || contains(t.ShowCompany, res) {
				perm.Actions = append(perm.Actions, "show")
			}
			if contains(t.EditCompany, assoicRes) || contains(t.EditCompany, res) {
				perm.Actions = append(perm.Actions, "edit")
			}
			if contains(t.DeleteCompany, assoicRes) || contains(t.DeleteCompany, res) {
				perm.Actions = append(perm.Actions, "delete")
			}

			perms = append(perms, bindPermID(&perm))
		}

		roles = append(roles, &Role{
			Model: resource.Model{
				ID: uint(len(roles) + 1),
			},
			Name:        t.Name,
			Type:        t.Type,
			Permissions: perms,
		})
	}
	return roles
}
