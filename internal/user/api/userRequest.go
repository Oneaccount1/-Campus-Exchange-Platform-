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
	Role      string    `json:"role"`
}

// PasswordUpdate 密码更新请求对象
type PasswordUpdate struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
