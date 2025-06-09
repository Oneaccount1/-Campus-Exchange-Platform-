package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/utils/config"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/utils/errors"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/utils/response"
	"strings"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.HandleError(c, errors.NewUnauthorizedError("无效的认证格式", nil))
			c.Abort()
			return
		}

		// 获取token
		tokenString := parts[1]

		// 解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewUnauthorizedError("无效的令牌", nil)
			}

			// 返回签名密钥
			return []byte(config.GetJWTConfig().Secret), nil
		})

		if err != nil {
			response.HandleError(c, errors.NewUnauthorizedError("无效的令牌", err))
			c.Abort()
			return
		}

		// 验证token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户信息保存到上下文
			c.Set("user_id", uint(claims["user_id"].(float64)))
			c.Set("username", claims["username"].(string))
			c.Set("role", claims["role"].(string))
			c.Next()
		} else {
			response.HandleError(c, errors.NewUnauthorizedError("无效的令牌", nil))
			c.Abort()
			return
		}
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行JWT认证
		JWTAuth()(c)

		// 如果认证失败，中间件链会中断，不会执行后续代码
		if c.IsAborted() {
			return
		}

		// 检查用户角色
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			response.HandleError(c, errors.NewForbiddenError("需要管理员权限", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
