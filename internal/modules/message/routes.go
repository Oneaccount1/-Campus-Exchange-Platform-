package message

import (
	"campus/internal/bootstrap"
	"campus/internal/middleware"
	"campus/internal/modules/message/controllers"
	"campus/internal/modules/message/repositories"
	"campus/internal/modules/message/services"
	"campus/internal/rabbitMQ"
	"campus/internal/utils/errors"
	"campus/internal/utils/logger"
	"campus/internal/utils/response"
	"campus/internal/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterRoutes 注册消息模块的路由
func RegisterRoutes(r *gin.Engine, api *gin.RouterGroup, wsManager *websocket.Manager, rabbitMQURL string) {
	// --- Dependency Injection ---

	// 1. Create Repository
	db := bootstrap.GetDB()
	messageRepo := repositories.NewMessageRepository(db)

	// 2. Create RabbitMQ Publisher
	// Note: Error handling should be more robust in a real application (e.g., panic or graceful shutdown)
	publisher, err := rabbitMQ.NewPublisher(rabbitMQURL)
	if err != nil {
		logger.Fatal("Failed to create RabbitMQ publisher", zap.Error(err))
	}

	// 3. Create Service
	messageService := services.NewMessageService(messageRepo, publisher)

	// --- Controller and Routes Setup ---

	controller := controllers.NewMessageController(messageService)

	// Message related REST API routes - authentication required
	messageGroup := api.Group("/messages")
	messageGroup.Use(middleware.JWTAuth())
	{
		messageGroup.POST("", controller.SendMessage)
		messageGroup.GET("/contacts", controller.GetContacts)
		messageGroup.GET("/:contactId", controller.GetMessages)
		messageGroup.GET("/:contactId/last", controller.GetLastMessage)
		messageGroup.PUT("/:contactId/read", controller.MarkAsRead)
		messageGroup.GET("/unread/count", controller.GetUnreadCount)
		messageGroup.POST("/conversation", controller.CreateConversation)
	}

	// WebSocket route - uses a dedicated WebSocket authentication middleware
	wsRoute := api.Group("/messages")
	wsRoute.Use(middleware.WSAuth())
	{
		// WebSocket endpoint for real-time messaging
		wsRoute.GET("/ws", func(c *gin.Context) {
			// Extract user ID from the context (set by the WSAuth middleware)
			userID, exists := c.Get("user_id")
			if !exists {
				response.HandleError(c, errors.ErrUnauthorized)
				return
			}

			// Upgrade the HTTP connection to a WebSocket connection
			wsManager.HandleConnection(c.Writer, c.Request, userID.(uint))
		})
	}
	
	// 管理员消息路由 - 需要管理员权限
	adminMessageGroup := api.Group("/admin/messages")
	adminMessageGroup.Use(middleware.JWTAuth())
	adminMessageGroup.Use(middleware.AuthorizeByRole("admin"))
	{
		// 获取消息列表
		adminMessageGroup.GET("", middleware.AuthorizePermission("/api/v1/admin/messages", "GET"), controller.GetAdminMessageList)
		
		// 获取会话列表
		adminMessageGroup.GET("/conversations", middleware.AuthorizePermission("/api/v1/admin/messages/conversations", "GET"), controller.GetAdminConversationList)
		
		// 获取会话消息历史
		adminMessageGroup.GET("/history", middleware.AuthorizePermission("/api/v1/admin/messages/history", "GET"), controller.GetAdminMessageHistory)
		
		// 发送系统消息
		adminMessageGroup.POST("/system", middleware.AuthorizePermission("/api/v1/admin/messages/system", "POST"), controller.SendSystemMessage)
		
		// 删除消息
		adminMessageGroup.DELETE("/:messageId", middleware.AuthorizePermission("/api/v1/admin/messages/:messageId", "DELETE"), controller.DeleteMessage)
	}
}
