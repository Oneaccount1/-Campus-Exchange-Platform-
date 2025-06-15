package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ProductID uint   `gorm:"not null"`
	ImageURL  string `gorm:"not null"`
}
