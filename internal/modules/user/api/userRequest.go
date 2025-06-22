package api

import "time"

// UserLogin 用户登录请求对象
type UserLogin struct {
	UserName string `json:"user_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
}

// UserRegister 用户注册数据传输对象
type UserRegister struct {
	Username    string `json:"user_name" binding:"required,min=3,max=50"`
	Password    string `json:"pass_word" binding:"required,min=6,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Nickname    string `json:"nickname"`
	Phone       string `json:"phone"`
	Description string `json:"description"` // 个性签名
}

// UserUpdate 用户更新数据传输对象
type UserUpdate struct {
	Nickname    string `json:"nickname"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

// JWTResponse JWT响应
type JWTResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Roles     []string  `json:"roles"` // 用户所有角色
}

// PasswordUpdate 密码更新请求对象
type PasswordUpdate struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// AdminUserListQuery 管理员用户列表查询参数
type AdminUserListQuery struct {
	Page      int    `form:"page" json:"page"`            // 页码
	Size      int    `form:"size" json:"pageSize"`        // 每页数量
	Search    string `form:"search" json:"search"`        // 搜索关键词
	Status    string `form:"status" json:"status"`        // 状态筛选
	StartDate string `form:"start_date" json:"start_ate"` // 开始日期
	EndDate   string `form:"end_date" json:"end_ate"`     // 结束日期
}

// UserStatusUpdate 用户状态更新请求
type UserStatusUpdate struct {
	Status string `json:"status" binding:"required"` // 用户状态：正常、禁用
}
