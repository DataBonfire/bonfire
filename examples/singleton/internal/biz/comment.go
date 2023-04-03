package biz

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID  uint
	Content string

	OrganizationID uint `gorm:"index"`
}
