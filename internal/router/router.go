package router

import (
	"campus/internal/bootstrap"
	Message "campus/internal/modules/message"
	Order "campus/internal/modules/order"
	Permission "campus/internal/modules/permission"
	Product "campus/internal/modules/product"
	User "campus/internal/modules/user"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本
	api := r.Group("/api/v1")

	// 注册各个模块的路由
	registerModuleRoutes(r, api)
}

// registerModuleRoutes 注册各个模块的路由
func registerModuleRoutes(r *gin.Engine, api *gin.RouterGroup) {
	// 获取配置和WebSocket管理器
	config := bootstrap.GetConfig()
	wsManager := bootstrap.GetWebSocketManager()

	// 用户模块路由
	User.RegisterRoutes(r, api)

	// 订单模块路由
	Order.RegisterRoutes(r, api)

	// 商品模块路由
	Product.RegisterRoutes(r, api)

	// 权限模块路由
	Permission.RegisterRoutes(r, api)

	// 消息模块路由
	Message.RegisterRoutes(r, api, wsManager, config.RabbitMQ.URL)
}
