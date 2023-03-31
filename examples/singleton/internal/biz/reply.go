package biz

import "gorm.io/gorm"

type Reply struct {
	gorm.Model
	CommentID uint
	PostID    uint

	CreatedBy      uint `gorm:"index"`
	OrganizationID uint `gorm:"index"`
}
