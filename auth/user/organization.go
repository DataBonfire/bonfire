package user

import "github.com/databonfire/bonfire/resource"

var OrganizationResourceName = "organizations"

type Organization struct {
	resource.Model
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Address string `json:"address"`
}

func (o Organization) TableName() string {
	return OrganizationResourceName
}
