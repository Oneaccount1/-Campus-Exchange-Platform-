package middleware

import (
	"campus/internal/modules/permission/repositories"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Authorize 基于Casbin的权限验证中间件
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 获取请求的路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 将userID转换为字符串
		sub := fmt.Sprintf("%d", userID.(uint))

		// 创建权限仓库
		permissionRepo := repositories.NewPermissionRepository()

		// 检查权限
		ok, err := permissionRepo.Enforce(sub, obj, act)
		if err != nil {
			response.HandleError(c, errors.NewInternalServerError("权限检查失败", err))
			c.Abort()
			return
		}

		if !ok {
			response.HandleError(c, errors.NewForbiddenError("没有操作权限", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthorizeByRole 基于角色的权限验证中间件
func AuthorizeByRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 首先检查上下文中的角色列表
		if roleList, exists := c.Get("roles"); exists {
			if roles, ok := roleList.([]string); ok {
				for _, r := range roles {
					if r == role {
						// 找到匹配的角色，允许访问
						c.Next()
						return
					}
				}
			}
		}

		// 将userID转换为字符串
		sub := fmt.Sprintf("%d", userID.(uint))

		// 创建权限仓库
		permissionRepo := repositories.NewPermissionRepository()

		// 检查用户是否有指定角色
		ok, err := permissionRepo.HasRoleForUser(sub, role)
		if err != nil {
			response.HandleError(c, errors.NewInternalServerError("角色检查失败", err))
			c.Abort()
			return
		}

		if !ok {
			response.HandleError(c, errors.NewForbiddenError("需要"+role+"角色权限", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthorizePermission 基于特定权限的验证中间件
func AuthorizePermission(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 将userID转换为字符串
		sub := fmt.Sprintf("%d", userID.(uint))

		// 创建权限仓库
		permissionRepo := repositories.NewPermissionRepository()

		// 检查权限
		ok, err := permissionRepo.Enforce(sub, obj, act)
		if err != nil {
			response.HandleError(c, errors.NewInternalServerError("权限检查失败", err))
			c.Abort()
			return
		}

		if !ok {
			// 检查用户角色列表中是否有admin角色，admin拥有所有权限
			if roleList, exists := c.Get("roles"); exists {
				if roles, ok := roleList.([]string); ok {
					for _, r := range roles {
						if r == "admin" {
							// 管理员拥有所有权限
							c.Next()
							return
						}
					}
				}
			}

			response.HandleError(c, errors.NewForbiddenError("没有操作权限", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
