package User

import (
	"campus/internal/middleware"
	"campus/internal/modules/user/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册User模块的所有路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	userController := controllers.NewUserController()
	favoriteController := controllers.NewFavoriteController()

	// 公开路由 - 不需要认证
	publicGroup := api.Group("/")
	registerPublicRoutes(publicGroup, userController)

	// 用户路由 - 需要认证
	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWTAuth())
	registerAuthRoutes(userGroup, userController)

	// 收藏路由 - 需要认证
	favoriteGroup := api.Group("/user/favorites")
	favoriteGroup.Use(middleware.JWTAuth())
	registerFavoriteRoutes(favoriteGroup, favoriteController)

	// 管理员路由 - 需要管理员权限
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.JWTAuth())
	adminGroup.Use(middleware.AuthorizeByRole("admin"))
	registerAdminRoutes(adminGroup, userController)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.POST("/admin/login", controller.AdminLogin)

}

// registerAuthRoutes 注册需要认证的路由
func registerAuthRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	// 个人资料 - 使用基于特定权限的中间件
	profileGroup := router.Group("/profile")
	//profileGroup.GET("", middleware.AuthorizePermission("/api/v1/user/profile", "GET"), controller.GetProfile)
	//profileGroup.PUT("", middleware.AuthorizePermission("/api/v1/user/profile", "PUT"), controller.UpdateProfile)
	profileGroup.GET("", controller.GetProfile)
	profileGroup.PUT("", controller.UpdateProfile)

	// 修改密码 - 使用基于特定权限的中间件
	//router.POST("/change-password", middleware.AuthorizePermission("/api/v1/user/change-password", "POST"), controller.ChangePassword)
	router.POST("/change-password", controller.ChangePassword)

	// 查看用户信息 - 使用基于特定权限的中间件
	//router.GET("/:id", middleware.AuthorizePermission("/api/v1/user/:id", "GET"), controller.GetUserByID)
	router.GET("/:id", controller.GetUserByID)
}

// registerFavoriteRoutes 注册收藏相关的路由
func registerFavoriteRoutes(router *gin.RouterGroup, controller *controllers.FavoriteController) {
	router.POST("", controller.AddFavorite)                   // 添加收藏
	router.DELETE("/:productID", controller.RemoveFavorite)   // 取消收藏
	router.GET("", controller.ListFavorites)                  // 获取收藏列表
	router.GET("/check/:productID", controller.CheckFavorite) // 检查是否已收藏
}

// registerAdminRoutes 注册需要管理员权限的路由
func registerAdminRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	// 用户列表接口（支持基本和高级查询）
	router.GET("/users", middleware.AuthorizePermission("/api/v1/admin/users", "GET"), controller.ListUsers)

	// 获取用户详情
	router.GET("/users/:id", middleware.AuthorizePermission("/api/v1/admin/users/:id", "GET"), controller.GetUserDetail)

	// 更新用户状态
	router.PUT("/users/:id/status", middleware.AuthorizePermission("/api/v1/admin/users/:id/status", "PUT"), controller.UpdateUserStatus)

	// 重置用户密码
	router.POST("/users/:id/reset-password", middleware.AuthorizePermission("/api/v1/admin/users/:id/reset-password", "POST"), controller.ResetUserPassword)
}
