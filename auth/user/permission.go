package user

import (
	"github.com/databonfire/bonfire/resource"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Action   string // list,read,create,edit,delete,export,print
	Resource string // *, campaigns
	Record   resource.Filter
}

// 1. me 2. org 3. 下属
// 不考虑同事

// {"action:"edit", "resource": "campaigns", "record": {"created_by.manager_id": "me", "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"created_by": "me"}, "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"campaign_id": "mine"}, "status": {"lt": 3}}
// campaigns -> list -> show edit btn?
// hasPermission(action, resource)

func m()  {
	
}