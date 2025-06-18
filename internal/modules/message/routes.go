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
		wsRoute.GET("/ws", func(ctx *gin.Context) {
			userID, exists := ctx.Get("user_id")
			if !exists {
				logger.Warn("WebSocket连接失败, 用户ID不存在")
				response.HandleError(ctx, errors.ErrUnauthorized)
				return
			}

			userIDUint := userID.(uint)
			logger.Info("处理websocket连接请求", zap.Uint("UserID", userIDUint))

			// The WebSocket handler's responsibility is now only to manage the connection.
			// All message processing logic (offline, online push) is handled by the
			// backend consumer listening to RabbitMQ.
			wsManager.HandleConnection(ctx.Writer, ctx.Request, userIDUint)

			// The complex and inefficient polling goroutine has been removed.
			// The client is now expected to fetch message history via a separate HTTP request.
		})
	}
}
