package models

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	gorm.Model
	Resource    string `gorm:"size:100;not null" json:"resource"`                                             // 资源路径
	Action      string `gorm:"size:50;not null" json:"action"`                                                // 操作类型 (GET, POST, PUT, DELETE等)
	Description string `gorm:"size:200" json:"description"`                                                   // 权限描述
	Roles       []Role `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"roles,omitempty"` // 拥有此权限的角色
}
