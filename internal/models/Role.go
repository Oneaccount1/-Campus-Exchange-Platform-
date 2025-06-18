package models

import "gorm.io/gorm"

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string       `gorm:"size:50;not null;uniqueIndex" json:"name"`                                            // 角色名称
	Description string       `gorm:"size:200" json:"description"`                                                         // 角色描述
	Users       []User       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"users,omitempty"`             // 拥有此角色的用户
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"permissions,omitempty"` // 角色拥有的权限
}
