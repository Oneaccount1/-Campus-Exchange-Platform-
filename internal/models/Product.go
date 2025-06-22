package models

import (
	"gorm.io/gorm"
	"time"
)

// Product 商品模型
type Product struct {
	gorm.Model
	Title         string         `gorm:"size:100;not null;index" json:"title"`
	Description   string         `gorm:"size:1000" json:"description"`
	Price         float64        `gorm:"not null" json:"price"`
	ProductImages []ProductImage `gorm:"size:1000" json:"images"` // JSON string of image URLs
	Category      string         `gorm:"size:50;index" json:"category"`
	Condition     string         `gorm:"size:20" json:"condition"` // new, like_new, good, fair, poor
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	Status        string         `gorm:"size:20;default:available" json:"status"` // 售卖中, 已下架, 审核中
	SoldAt        time.Time      `json:"sold_at"`
}
