package User

import (
	"campus/internal/middleware"
	"campus/internal/modules/user/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册User模块的所有路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	userController := controllers.NewUserController()

	// 公开路由 - 不需要认证
	publicGroup := api.Group("/auth")
	registerPublicRoutes(publicGroup, userController)

	// 用户路由 - 需要认证
	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWTAuth())
	registerAuthRoutes(userGroup, userController)

	// 管理员路由 - 需要管理员权限
	adminGroup := api.Group("/admin/users")
	adminGroup.Use(middleware.JWTAuth())
	adminGroup.Use(middleware.AuthorizeByRole("admin"))
	registerAdminRoutes(adminGroup, userController)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
}

// registerAuthRoutes 注册需要认证的路由
func registerAuthRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	router.GET("/profile", controller.GetProfile)
	router.PUT("/profile", controller.UpdateProfile)
	router.POST("/change-password", controller.ChangePassword)
	router.GET("/:id", controller.GetUserByID)
}

// registerAdminRoutes 注册需要管理员权限的路由
func registerAdminRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	router.GET("", controller.ListUsers)
}
