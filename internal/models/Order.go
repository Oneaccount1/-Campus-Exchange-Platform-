package models

import (
	"gorm.io/gorm"
	"time"
)

// Order 订单模型
type Order struct {
	gorm.Model
	BuyerID      uint       `gorm:"not null;index" json:"buyer_id"`
	Buyer        User       `gorm:"foreignKey:BuyerID" json:"buyer"`
	SellerID     uint       `gorm:"not null;index" json:"seller_id"`
	Seller       User       `gorm:"foreignKey:SellerID" json:"seller"`
	ProductID    uint       `gorm:"not null;index" json:"product_id"`
	Product      Product    `gorm:"foreignKey:ProductID" json:"product"`
	Status       string     `gorm:"size:20;default:卖家未处理" json:"status"` // pending, completed, cancelled
	PayTime      *time.Time `json:"pay_time"`
	DeliveryTime *time.Time `json:"delivery_time"`
	CompleteTime *time.Time `json:"complete_time"`
	Remark       string     `json:"remark"`
}

// OrderLog 订单日志模型
type OrderLog struct {
	gorm.Model
	OrderID  uint   `json:"order_id"`
	Action   string `json:"action"`
	Operator string `json:"operator"`
	Remark   string `json:"remark"`
}
