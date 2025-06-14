package controllers

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"campus/internal/modules/permission/api"
	"campus/internal/modules/permission/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// PermissionController 权限控制器
type PermissionController struct {
	permissionService services.PermissionService
}

// NewPermissionController 创建权限控制器实例
func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: services.NewPermissionService(),
	}
}

// AssignRole 为用户分配角色
func (c *PermissionController) AssignRole(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 绑定请求
	var req api.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 开启事务
	tx := bootstrap.GetDB().Begin()
	if tx.Error != nil {
		response.HandleError(ctx, errors.NewInternalServerError("开启事务失败", tx.Error))
		return
	}

	// 查找用户
	var user models.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewNotFoundError("用户", err))
		return
	}

	// 查找角色
	var role models.Role
	if err := tx.Where("name = ?", req.Role).First(&role).Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewNotFoundError("角色", err))
		return
	}

	// 检查是否已经有该角色
	var count int64
	if err := tx.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", user.ID, role.ID).Count(&count).Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewInternalServerError("检查角色关联失败", err))
		return
	}

	if count == 0 {
		// 添加角色关联
		if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
			tx.Rollback()
			response.HandleError(ctx, errors.NewInternalServerError("关联角色失败", err))
			return
		}
	}

	// 分配角色到Casbin
	err = c.permissionService.AddRoleForUser(uint(id), req.Role)
	if err != nil {
		tx.Rollback()
		response.HandleError(ctx, err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewInternalServerError("提交事务失败", err))
		return
	}

	response.SuccessWithMessage(ctx, "角色分配成功", nil)
}

// RemoveRole 移除用户角色
func (c *PermissionController) RemoveRole(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 绑定请求
	var req api.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 开启事务
	tx := bootstrap.GetDB().Begin()
	if tx.Error != nil {
		response.HandleError(ctx, errors.NewInternalServerError("开启事务失败", tx.Error))
		return
	}

	// 查找用户
	var user models.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewNotFoundError("用户", err))
		return
	}

	// 查找角色
	var role models.Role
	if err := tx.Where("name = ?", req.Role).First(&role).Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewNotFoundError("角色", err))
		return
	}

	// 移除角色关联
	if err := tx.Model(&user).Association("Roles").Delete(&role); err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewInternalServerError("移除角色关联失败", err))
		return
	}

	// 从Casbin移除角色
	err = c.permissionService.DeleteRoleForUser(uint(id), req.Role)
	if err != nil {
		tx.Rollback()
		response.HandleError(ctx, err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.HandleError(ctx, errors.NewInternalServerError("提交事务失败", err))
		return
	}

	response.SuccessWithMessage(ctx, "角色移除成功", nil)
}

// GetUserRoles 获取用户角色
func (c *PermissionController) GetUserRoles(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 获取角色
	roles, err := c.permissionService.GetRolesForUser(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, &api.RoleListResponse{Roles: roles})
}

// AddPermission 添加权限
func (c *PermissionController) AddPermission(ctx *gin.Context) {
	var req api.PermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	err := c.permissionService.AddPermissionForRole(req.Role, req.Object, req.Action)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "权限添加成功", nil)
}

// RemovePermission 移除权限
func (c *PermissionController) RemovePermission(ctx *gin.Context) {
	var req api.PermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	err := c.permissionService.DeletePermissionForRole(req.Role, req.Object, req.Action)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "权限移除成功", nil)
}

// CheckPermission 检查权限
func (c *PermissionController) CheckPermission(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req api.CheckPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 检查权限
	hasPermission, err := c.permissionService.CheckPermission(userID.(uint), req.Object, req.Action)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{"has_permission": hasPermission})
}

// GetUserPermissions 获取用户权限
func (c *PermissionController) GetUserPermissions(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效用户ID", err))
		return
	}

	// 获取权限
	permissions, err := c.permissionService.GetUserPermissions(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	// 添加用户ID到响应中
	response.Success(ctx, gin.H{
		"user_id":     id,
		"permissions": permissions.Permissions,
		"roles": func() []string {
			roles, _ := c.permissionService.GetRolesForUser(uint(id))
			return roles
		}(),
	})
}
