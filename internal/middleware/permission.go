package middleware

import (
	"campus/internal/bootstrap"
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

		// 检查权限
		enforcer := bootstrap.GetEnforcer()
		ok, err := enforcer.Enforce(sub, obj, act)
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

		// 将userID转换为字符串
		sub := fmt.Sprintf("%d", userID.(uint))

		// 检查用户是否有指定角色
		enforcer := bootstrap.GetEnforcer()
		ok, err := enforcer.HasRoleForUser(sub, role)
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

		// 检查权限
		enforcer := bootstrap.GetEnforcer()
		ok, err := enforcer.Enforce(sub, obj, act)
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
