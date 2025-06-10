package router

import (
	User "campus/internal/modules/user"
	"github.com/gin-gonic/gin"
)

// 所有模块路由在这里注册

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本
	api := r.Group("/api/v1")

	// 注册各个模块的路由
	registerModuleRoutes(r, api)
}

// registerModuleRoutes 注册各个模块的路由
func registerModuleRoutes(r *gin.Engine, api *gin.RouterGroup) {
	// 用户模块路由
	User.RegisterRoutes(r, api)

}
