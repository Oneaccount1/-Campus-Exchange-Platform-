package models

import "gorm.io/gorm"

// Message 消息模型
type Message struct {
	gorm.Model
	SenderID   uint   `gorm:"not null;index" json:"sender_id"`
	Sender     User   `gorm:"foreignKey:SenderID" json:"sender"`
	ReceiverID uint   `gorm:"not null;index" json:"receiver_id"`
	Receiver   User   `gorm:"foreignKey:ReceiverID" json:"receiver"`
	Content    string `gorm:"size:1000;not null" json:"content"`
	IsRead     bool   `gorm:"default:false" json:"is_read"`
	ProductID  uint   `gorm:"index" json:"product_id"`
	//Product    Product `gorm:"foreignKey:ProductID" json:"product"`
}
