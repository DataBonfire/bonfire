package biz

import "gorm.io/gorm"

type Influencer struct {
	gorm.Model

	Name                    string
	UserId                  uint
	Gender                  string
	Avatar                  string
	Yob                     string
	Identity                string
	Country                 string
	State                   string
	City                    string
	Address                 string
	EmailForCooperation     string
	TelephoneForCooperation string
	Biography               string
	IsVerified              bool

	CreatedBy      uint `gorm:"index"`
	OrganizationID uint `gorm:"index"`
}
