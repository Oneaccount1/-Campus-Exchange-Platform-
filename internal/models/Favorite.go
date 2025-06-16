package models

import (
	"gorm.io/gorm"
)

// Favorite 收藏模型
type Favorite struct {
	gorm.Model
	UserID    uint    `gorm:"not null;index" json:"user_id"`       // 用户ID
	ProductID uint    `gorm:"not null;index" json:"product_id"`    // 商品ID
	Product   Product `gorm:"foreignKey:ProductID" json:"product"` // 关联的商品
}
