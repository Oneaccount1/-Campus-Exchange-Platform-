package upload

import (
	"campus/internal/middleware"
	"campus/internal/modules/upload/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册上传相关的路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup) {
	// 创建控制器实例
	uploadController := controllers.NewUploadController()
	
	// 上传路由组 - 需要认证
	uploadGroup := api.Group("/upload")
	uploadGroup.Use(middleware.JWTAuth())
	
	// 上传图片
	uploadGroup.POST("/image", uploadController.UploadImage)
	
	// 上传通用文件
	uploadGroup.POST("/file", uploadController.UploadFile)
	
	// 为产品模块添加专用上传路由
	productUploadGroup := api.Group("/product")
	productUploadGroup.Use(middleware.JWTAuth())
	productUploadGroup.POST("/upload", uploadController.UploadImage)
	
	// 为用户模块添加专用上传路由
	userUploadGroup := api.Group("/user")
	userUploadGroup.Use(middleware.JWTAuth())
	userUploadGroup.POST("/avatar", uploadController.UploadImage)
	
	// 为消息模块添加专用上传路由
	messageUploadGroup := api.Group("/messages")
	messageUploadGroup.Use(middleware.JWTAuth())
	messageUploadGroup.POST("/upload", uploadController.UploadImage)
	
	// 添加静态文件服务
	r.Static("/static", "./uploads")
} 