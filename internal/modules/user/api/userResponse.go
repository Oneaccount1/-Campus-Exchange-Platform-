package api

import "time"

// UserResponse 用户信息响应
type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Avatar       string    `json:"avatar"`
	Roles        []string  `json:"roles"` // 用户所有角色
	Description  string    `json:"description"`
	Status       string    `json:"status"` // 用户状态
	ProductCount int       `json:"product_count"` // 用户发布的产品数量
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total"`
	List  []UserResponse `json:"list"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// AdminUserListResponse 管理员用户列表响应
type AdminUserListResponse struct {
	Total int64                `json:"total"`
	List  []AdminUserResponse  `json:"list"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// AdminUserResponse 管理员用户信息响应
type AdminUserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Avatar       string    `json:"avatar"`
	RegisterTime time.Time `json:"registerTime"` // 注册时间
	ProductCount int       `json:"productCount"` // 产品数量
	Status       string    `json:"status"`       // 用户状态
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	NewPassword string `json:"newPassword"` // 新密码
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	ID            uint                  `json:"id"`
	Username      string                `json:"username"`
	Avatar        string                `json:"avatar"`
	RegisterTime  time.Time             `json:"registerTime"`  // 注册时间
	Email         string                `json:"email"`
	Phone         string                `json:"phone"`
	LastLogin     time.Time             `json:"lastLogin"`     // 最后登录时间
	LastIP        string                `json:"lastIp"`        // 最后登录IP
	Status        string                `json:"status"`        // 用户状态
	ProductCount  int                   `json:"productCount"`  // 产品数量
	OrderCount    int                   `json:"orderCount"`    // 订单数量
	FavoriteCount int                   `json:"favoriteCount"` // 收藏数量
	Products      []UserProductItem     `json:"products"`      // 用户发布的商品
	Activities    []UserActivityItem    `json:"activities"`    // 用户活动
}

// UserProductItem 用户商品项
type UserProductItem struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Price      float64   `json:"price"`
	CreateTime time.Time `json:"createTime"`
}

// UserActivityItem 用户活动项
type UserActivityItem struct {
	Content string    `json:"content"` // 活动内容
	Time    time.Time `json:"time"`    // 活动时间
}
