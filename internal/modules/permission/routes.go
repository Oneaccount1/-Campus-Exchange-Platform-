package permission

import (
	"campus/internal/middleware"
	"campus/internal/modules/permission/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册权限模块路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	permissionController := controllers.NewPermissionController()

	// 权限管理路由 -需要管理员权限
	permissionGroup := api.Group("admin/permissions")
	permissionGroup.Use(middleware.JWTAuth())
	permissionGroup.Use(middleware.AuthorizeByRole("admin"))

	// 角色管理
	permissionGroup.POST("/users/:id/roles", permissionController.AssignRole)
	permissionGroup.DELETE("/users/:id/roles", permissionController.RemoveRole)
	permissionGroup.GET("/users/:id/roles", permissionController.GetUserRoles)

	// 权限管理
	permissionGroup.POST("/policies", permissionController.AddPermission)
	permissionGroup.DELETE("/policies", permissionController.RemovePermission)

	// 权限检查
	permissionGroup.POST("/check", permissionController.CheckPermission)
	permissionGroup.GET("/users/:id/permissions", permissionController.GetUserPermissions)
}
