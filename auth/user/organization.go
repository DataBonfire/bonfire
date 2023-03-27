package user

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	Name    string
	Logo    string
	Address string
}
