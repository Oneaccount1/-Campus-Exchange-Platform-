package models

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Username    string `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password    string `gorm:"size:100;not null" json:"-"`
	Nickname    string `gorm:"size:50" json:"nickname"`
	Email       string `gorm:"size:100;uniqueIndex" json:"email"`
	Phone       string `gorm:"size:20" json:"phone"`
	Avatar      string `gorm:"size:255" json:"avatar"`
	Roles       []Role `gorm:"many2many:user_roles" json:"roles,omitempty"` // 用户拥有的所有角色
	Description string `gorm:"size:500" json:"description"`
}
