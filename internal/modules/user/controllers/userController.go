package controllers

import (
	"campus/internal/modules/user/api"
	"campus/internal/modules/user/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"

	"strconv"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req api.UserRegister
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	user, err := c.userService.Register(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "注册成功", user)
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req api.UserLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	// 普通用户登录，不需要验证特定角色
	token, err := c.userService.Login(&req, "")
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	// 成功返回token
	response.Success(ctx, token)
}

// AdminLogin 管理员登录
func (c *UserController) AdminLogin(ctx *gin.Context) {
	// 从请求中获取用户名和密码
	var req api.UserLogin

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	// 调用登录服务，验证管理员角色
	token, err := c.userService.Login(&req, "admin")
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	// 成功返回token
	response.Success(ctx, token)
}

// GetProfile 获取用户个人资料
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}
	// 根据ID查用户
	user, err := c.userService.GetByID(userID.(uint))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, user)
}

// UpdateProfile 更新用户个人资料
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req api.UserUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	user, err := c.userService.UpdateUser(userID.(uint), req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "更新成功", user)
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req api.PasswordUpdate

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	err := c.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "密码修改成功", nil)
}

// GetUserByID 根据ID获取用户信息
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, user)
}

// ListUsers 获取用户列表（管理员功能）
func (c *UserController) ListUsers(ctx *gin.Context) {
	// 解析所有查询参数
	var query api.AdminUserListQuery

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 如果没有使用高级查询参数，使用基本列表功能
	if query.Search == "" && query.Status == "" && query.StartDate == "" && query.EndDate == "" {
		// 设置默认值
		if query.Page <= 0 {
			query.Page = 1
		}
		if query.Size <= 0 {
			query.Size = 10
		}

		result, err := c.userService.List(query.Page, query.Size)
		if err != nil {
			response.HandleError(ctx, err)
			return
		}
		response.Success(ctx, result)
		return
	}

	// 使用高级查询
	result, err := c.userService.AdminList(&query)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, result)
}

// AdminListUsers 管理员高级用户列表查询
func (c *UserController) AdminListUsers(ctx *gin.Context) {
	var query api.AdminUserListQuery

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	// 调用服务层方法
	result, err := c.userService.AdminList(&query)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "获取成功", result)
}

// UpdateUserStatus 更新用户状态
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 绑定请求参数
	var req api.UserStatusUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 调用服务层方法
	err = c.userService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "状态更新成功", nil)
}

// ResetUserPassword 重置用户密码
func (c *UserController) ResetUserPassword(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 调用服务层方法
	result, err := c.userService.ResetPassword(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "密码已重置", result)
}

// GetUserDetail 获取用户详情
func (c *UserController) GetUserDetail(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 调用服务层方法
	detail, err := c.userService.GetUserDetail(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "获取成功", detail)
}
