package biz

import "gorm.io/gorm"

type Reply struct {
	gorm.Model
	CommentID uint
	PostID    uint

	OrganizationID uint `gorm:"index"`
}
