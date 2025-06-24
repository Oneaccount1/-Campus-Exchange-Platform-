package product

import (
	"campus/internal/middleware"
	"campus/internal/modules/product/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册product模块的所有路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	productController := controllers.NewProductController()

	// 商品路由 - 需要认证
	productGroup := api.Group("/product")
	productGroup.Use(middleware.JWTAuth())
	registerProductRoutes(productGroup, productController)

	// 管理员商品路由 - 需要管理员权限
	adminProductGroup := api.Group("/admin/products")
	adminProductGroup.Use(middleware.JWTAuth())
	adminProductGroup.Use(middleware.AuthorizeByRole("admin"))
	registerAdminProductRoutes(adminProductGroup, productController)
}

// registerProductRoutes 注册商品相关路由
func registerProductRoutes(router *gin.RouterGroup, controller *controllers.ProductController) {
	router.GET("", controller.ListProducts)
	router.POST("", controller.CreateProduct)
	router.GET("/:id", controller.GetProductByID)
	router.PUT("/:id", controller.UpdateProduct)
	router.DELETE("/:id", controller.DeleteProduct)
	router.GET("/search", controller.SearchProductsByKeyword)
	router.GET("/user", controller.GetUserProducts)
	router.GET("/solving", controller.ListSolvingProducts)
	router.GET("/latest", controller.GetLatestProducts)

}

// registerAdminProductRoutes 注册管理员商品相关路由
func registerAdminProductRoutes(router *gin.RouterGroup, controller *controllers.ProductController) {
	// 获取商品列表（基本功能和前台相同）
	router.GET("", middleware.AuthorizePermission("/api/v1/admin/products", "GET"), controller.AdminListProducts)

	// 获取商品详情（基本功能和前台相同）
	router.GET("/:id", middleware.AuthorizePermission("/api/v1/admin/products/:id", "GET"), controller.GetProductByID)

	// 更新商品（基本功能和前台相同）
	router.PUT("/:id", middleware.AuthorizePermission("/api/v1/admin/products/:id", "PUT"), controller.UpdateProduct)

	// 删除商品（基本功能和前台相同）
	router.DELETE("/:id", middleware.AuthorizePermission("/api/v1/admin/products/:id", "DELETE"), controller.DeleteProduct)

	// 更新商品状态
	router.PUT("/:id/status", middleware.AuthorizePermission("/api/v1/admin/products/:id/status", "PUT"), controller.UpdateProductStatus)

	// 批量更新商品状态
	//router.PUT("/batch-status", middleware.AuthorizePermission("/api/v1/admin/products/batch-status", "PUT"), controller.BatchUpdateStatus)

	// 最新商品（仪表盘使用的功能）
	//router.GET("/latest", middleware.AuthorizePermission("/api/v1/admin/products/latest", "GET"), controller.GetLatestProducts)
}
