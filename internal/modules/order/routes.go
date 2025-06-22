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
	
	// 管理员订单路由 - 需要管理员权限
	adminOrderGroup := api.Group("/admin/orders")
	adminOrderGroup.Use(middleware.JWTAuth())
	adminOrderGroup.Use(middleware.AuthorizeByRole("admin"))
	registerAdminOrderRoutes(adminOrderGroup, orderController)
}

// registerOrderRoutes 注册订单相关路由
func registerOrderRoutes(router *gin.RouterGroup, controller *controllers.OrderController) {
	router.POST("", controller.CreateOrder)
	router.DELETE("/:id", controller.DeleteOrder)
	router.PUT("/:id/status", controller.UpdateOrderStatus)
	router.GET("/:id", controller.GetOrderByID)
	router.GET("/user", controller.GetUserOrders)
}

// registerAdminOrderRoutes 注册管理员订单相关路由
func registerAdminOrderRoutes(router *gin.RouterGroup, controller *controllers.OrderController) {
	// 获取订单列表
	router.GET("", middleware.AuthorizePermission("/api/v1/admin/orders", "GET"), controller.GetAdminOrderList)
	
	// 获取订单详情
	router.GET("/:id", middleware.AuthorizePermission("/api/v1/admin/orders/:id", "GET"), controller.GetAdminOrderDetail)
	
	// 更新订单状态
	router.PUT("/:id/status", middleware.AuthorizePermission("/api/v1/admin/orders/:id/status", "PUT"), controller.AdminUpdateOrderStatus)
	
	// 导出订单数据
	router.GET("/export", middleware.AuthorizePermission("/api/v1/admin/orders/export", "GET"), controller.ExportOrders)
}
