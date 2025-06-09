package models

import (
	"gorm.io/gorm"
)

// Review 评价模型
type Review struct {
	gorm.Model
	UserID     uint    `gorm:"not null;index" json:"user_id"`
	User       User    `gorm:"foreignKey:UserID" json:"user"`
	ProductID  uint    `gorm:"not null;index" json:"product_id"`
	Product    Product `gorm:"foreignKey:ProductID" json:"product"`
	Rating     int     `gorm:"not null" json:"rating"` // 1-5
	Content    string  `gorm:"size:500" json:"content"`
	ReviewerID uint    `gorm:"not null;index" json:"reviewer_id"`
	Reviewer   User    `gorm:"foreignKey:ReviewerID" json:"reviewer"`
}
