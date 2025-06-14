package message

import (
	"campus/internal/middleware"
	"campus/internal/modules/message/controllers"
	"campus/internal/modules/message/services"
	"campus/internal/websocket"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册消息模块的路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup, wsManager *websocket.Manager, rabbitMQURL string) {
	// 创建消息服务
	messageService := services.NewMessageService(wsManager, rabbitMQURL)

	// 创建消息控制器
	controller := controllers.NewMessageController(messageService)

	// 消息相关路由 - 需要认证
	messageGroup := api.Group("/messages")
	messageGroup.Use(middleware.JWTAuth())
	{
		// REST API路由
		messageGroup.POST("", controller.SendMessage)               // 发送消息
		messageGroup.GET("/contacts", controller.GetContacts)       // 获取联系人列表
		messageGroup.GET("/:contactId", controller.GetMessages)     // 获取与联系人的消息历史
		messageGroup.PUT("/:contactId/read", controller.MarkAsRead) // 标记消息为已读

		// WebSocket路由
		messageGroup.GET("/ws", func(ctx *gin.Context) {
			// 获取当前用户ID
			userID, exists := ctx.Get("user_id")
			if !exists {
				ctx.JSON(401, gin.H{"error": "未授权"})
				return
			}

			// 处理WebSocket连接
			wsManager.HandleConnection(ctx.Writer, ctx.Request, userID.(uint))

			// 处理离线消息
			go messageService.ProcessOfflineMessages(userID.(uint))
		})
	}
}
