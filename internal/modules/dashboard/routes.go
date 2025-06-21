package dashboard

import (
	"campus/internal/middleware"
	"campus/internal/modules/dashboard/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册仪表盘模块的路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	dashboardController := controllers.NewDashboardController()

	// 仪表盘路由 - 需要管理员权限
	dashboardGroup := api.Group("/admin/dashboard")
	dashboardGroup.Use(middleware.JWTAuth())
	dashboardGroup.Use(middleware.AuthorizeByRole("admin"))

	// 统计数据
	dashboardGroup.GET("/stats", middleware.AuthorizePermission("/api/v1/admin/dashboard/stats", "GET"), dashboardController.GetStats)
	
	// 商品发布趋势
	dashboardGroup.GET("/product-trend", middleware.AuthorizePermission("/api/v1/admin/dashboard/product-trend", "GET"), dashboardController.GetProductTrend)
	
	// 商品分类统计
	dashboardGroup.GET("/category-stats", middleware.AuthorizePermission("/api/v1/admin/dashboard/category-stats", "GET"), dashboardController.GetCategoryStats)
	
	// 最新商品
	dashboardGroup.GET("/latest-products", middleware.AuthorizePermission("/api/v1/admin/dashboard/latest-products", "GET"), dashboardController.GetLatestProducts)
	
	// 系统活动
	dashboardGroup.GET("/activities", middleware.AuthorizePermission("/api/v1/admin/dashboard/activities", "GET"), dashboardController.GetActivities)
} 