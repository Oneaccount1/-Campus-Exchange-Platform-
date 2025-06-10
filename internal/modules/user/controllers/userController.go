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
	// 返回token 和 err
	token, err := c.userService.Login(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	// 成功者返回token
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
	// 查询URL中page参数如果没值就使用默认 1
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	result, err := c.userService.List(page, pageSize)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.Success(ctx, result)
}
