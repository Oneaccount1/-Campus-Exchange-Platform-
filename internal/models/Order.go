package models

import "gorm.io/gorm"

// Order 订单模型
type Order struct {
	gorm.Model
	BuyerID   uint    `gorm:"not null;index" json:"buyer_id"`
	Buyer     User    `gorm:"foreignKey:BuyerID" json:"buyer"`
	SellerID  uint    `gorm:"not null;index" json:"seller_id"`
	Seller    User    `gorm:"foreignKey:SellerID" json:"seller"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Status    string  `gorm:"size:20;default:pending" json:"status"` // pending, completed, cancelled
	Price     float64 `gorm:"not null" json:"price"`
}
