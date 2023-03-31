package biz

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title              string    `validate:"required" json:"title"`
	Content            string    `validate:"required" json:"content"`
	Status             int       `json:"status" user:"-"`
	PublishedAt        time.Time `json:"status"`
	CreatedBy          uint      `gorm:"index"`
	AdvertisingRevenue float64   `json:"-" admin:"advertising_revenue" editor_manager:"advertising_revenue" gorm:"advertising_revenue"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.PublishedAt = time.Now()
	return nil
}
