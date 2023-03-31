package biz

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string `validate:"required" json:"title"`
	Content     string `validate:"required" json:"content"`
	Status      int
	PublishedAt time.Time
	PublishedBy uint
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.PublishedAt = time.Now()
	return nil
}
