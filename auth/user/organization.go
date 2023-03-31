package user

import "github.com/databonfire/bonfire/resource"

type Organization struct {
	resource.Model
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Address string `json:"address"`
}
