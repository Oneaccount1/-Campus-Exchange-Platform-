package api

import "time"

// UserResponse 用户信息响应
type UserResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`      // 主要角色（向后兼容）
	RoleList    []string  `json:"role_list"` // 用户所有角色列表
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}
