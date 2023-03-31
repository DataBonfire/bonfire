package rbac

import "github.com/databonfire/bonfire/ac"

type Accessor interface {
	ac.Accessor
	GetRoles() []string
	GetRoleType() string
}

type rolevisitor string

var visitor rolevisitor = "visitor"

func (v *rolevisitor) GetID() uint {
	return 0
}

func (v *rolevisitor) GetGroups() []uint {
	return nil
}

func (v *rolevisitor) GetSubordinates() []uint {
	return nil
}

func (v *rolevisitor) GetRoles() []string {
	return []string{string(*v)}
}

func (v *rolevisitor) GetRoleType() string {
	return string(*v)
}
