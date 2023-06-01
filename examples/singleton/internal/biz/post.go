package biz

import (
	"time"

	"github.com/databonfire/bonfire/resource"
	"gorm.io/gorm"
)

type Post struct {
	resource.Model
	Title              string    `validate:"required" json:"title" editor:"editor_title" admin:"admin_title"`
	Content            string    `validate:"required" json:"content"`
	Status             int       `json:"status" user:"-"`
	PublishedAt        time.Time `json:"status"`
	AdvertisingRevenue float64   `json:"-" admin:"advertising_revenue" editor_manager:"advertising_revenue" gorm:"advertising_revenue"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.PublishedAt = time.Now()
	return nil
}
