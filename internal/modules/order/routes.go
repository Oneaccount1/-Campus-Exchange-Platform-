package order

import (
	"campus/internal/bootstrap"
	"campus/internal/middleware"
	"campus/internal/modules/order/controllers"
	"campus/internal/modules/order/repositories"
	"campus/internal/modules/order/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册order模块的所有路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	orderController := controllers.NewOrderController(services.NewOrderService(repositories.NewOrderRepository(bootstrap.GetDB())))

	// 订单路由 - 需要认证
	orderGroup := api.Group("/order")
	orderGroup.Use(middleware.JWTAuth())
	registerOrderRoutes(orderGroup, orderController)
}

// registerOrderRoutes 注册订单相关路由
func registerOrderRoutes(router *gin.RouterGroup, controller *controllers.OrderController) {
	router.POST("", controller.CreateOrder)
	router.DELETE("/:id", controller.DeleteOrder)
	router.PUT("/:id/status", controller.UpdateOrderStatus)
	router.GET("/:id", controller.GetOrderByID)
}
