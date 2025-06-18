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
}
