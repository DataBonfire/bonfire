package rbac

import (
	"github.com/databonfire/bonfire/filter"
	"github.com/databonfire/bonfire/resource"
)

type Permission struct {
	resource.Model
	Actions  resource.StringSlice `json:"actions"`  // list,read,create,edit,delete,export,print
	Resource string               `json:"resource"` // *, campaigns
	Record   filter.Filter        `json:"record"`
}

// 1. me 2. org 3. 下属
// 不考虑同事

// {"action:"edit", "resource": "campaigns", "record": {"created_by.manager_id": "me", "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"created_by": "me"}, "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"campaign_id": "mine"}, "status": {"lt": 3}}
// campaigns -> list -> show edit btn?
// hasPermission(action, resource)
