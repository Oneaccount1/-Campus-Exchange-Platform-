package message

import (
	"campus/internal/middleware"
	"campus/internal/modules/message/controllers"
	"campus/internal/modules/message/services"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		messageGroup.POST("", controller.SendMessage)                     // 发送消息
		messageGroup.GET("/contacts", controller.GetContacts)             // 获取联系人列表
		messageGroup.GET("/:contactId", controller.GetMessages)           // 获取与联系人的消息历史
		messageGroup.GET("/:contactId/last", controller.GetLastMessage)   // 获取与联系人的最后一条消息
		messageGroup.PUT("/:contactId/read", controller.MarkAsRead)       // 标记消息为已读
		messageGroup.GET("/unread/count", controller.GetUnreadCount)      // 获取未读消息数量
		messageGroup.POST("/conversation", controller.CreateConversation) // 创建新会话
	}

	// WebSocket路由 - 使用专门的WebSocket认证中间件
	wsRoute := api.Group("/messages")
	wsRoute.Use(middleware.WSAuth())
	{
		wsRoute.GET("/ws", func(ctx *gin.Context) {
			// 获取当前用户ID
			userID, exists := ctx.Get("user_id")
			if !exists {
				logger.Warn("WebSocket连接认证失败：未找到用户ID")
				ctx.JSON(401, gin.H{"error": "未授权"})
				return
			}

			userIDUint := userID.(uint)
			logger.Info("处理WebSocket连接请求", zap.Uint("用户ID", userIDUint))

			// 处理WebSocket连接
			wsManager.HandleConnection(ctx.Writer, ctx.Request, userIDUint)

			// 处理离线消息
			// 这里使用goroutine是安全的，因为HandleConnection中已经注册了连接
			go func(uid uint) {
				if err := messageService.ProcessOfflineMessages(uid); err != nil {
					logger.Error("处理离线消息失败", zap.Uint("用户ID", uid), zap.Error(err))
				}
			}(userIDUint)
		})
	}
}
