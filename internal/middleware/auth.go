package middleware

import (
	"campus/internal/bootstrap"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

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
			return []byte(bootstrap.GetConfig().JWT.Secret), nil
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

			// 处理角色列表
			if rolesRaw, ok := claims["roles"]; ok {
				if rolesArray, ok := rolesRaw.([]interface{}); ok {
					roles := make([]string, len(rolesArray))
					for i, r := range rolesArray {
						if roleStr, ok := r.(string); ok {
							roles[i] = roleStr
						}
					}
					c.Set("roles", roles)
				}
			}

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
	return AuthorizeByRole("admin")
}

// WSAuth WebSocket认证中间件
// 从URL查询参数获取token
func WSAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从URL查询参数获取token
		tokenString := c.Query("token")
		if tokenString == "" {
			response.HandleError(c, errors.NewUnauthorizedError("缺少token参数", nil))
			c.Abort()
			return
		}

		// 解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewUnauthorizedError("无效的令牌", nil)
			}

			// 返回签名密钥
			return []byte(bootstrap.GetConfig().JWT.Secret), nil
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

			// 处理角色列表
			if rolesRaw, ok := claims["roles"]; ok {
				if rolesArray, ok := rolesRaw.([]interface{}); ok {
					roles := make([]string, len(rolesArray))
					for i, r := range rolesArray {
						if roleStr, ok := r.(string); ok {
							roles[i] = roleStr
						}
					}
					c.Set("roles", roles)
				}
			}

			c.Next()
		} else {
			response.HandleError(c, errors.NewUnauthorizedError("无效的令牌", nil))
			c.Abort()
			return
		}
	}
}
