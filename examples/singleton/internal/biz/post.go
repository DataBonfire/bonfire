package biz

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string
	Content     string
	Status      int
	PublishedAt time.Time
	PublishedBy uint
}
