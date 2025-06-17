package message

import (
	"campus/internal/middleware"
	"campus/internal/modules/message/controllers"
	"campus/internal/modules/message/services"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
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
			// 使用goroutine并设置延迟，确保连接完全建立后再处理离线消息
			go func(uid uint) {
				// 给WebSocket连接一点时间建立
				time.Sleep(500 * time.Millisecond)

				// 记录处理开始
				logger.Info("开始处理用户离线消息", zap.Uint("用户ID", uid))

				// 添加panic恢复机制，防止goroutine崩溃
				defer func() {
					if r := recover(); r != nil {
						logger.Error("消息处理协程崩溃", zap.Uint("用户ID", uid), zap.Any("panic", r))
					}
				}()

				// 首先处理离线消息
				if err := messageService.ProcessOfflineMessages(uid); err != nil {
					logger.Error("处理离线消息失败", zap.Uint("用户ID", uid), zap.Error(err))
				}

				// 创建停止通道，用于安全退出
				stopChan := make(chan struct{})

				// 设置定期检查新消息的定时器
				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()

				// 启动一个监控goroutine检查用户状态
				go func() {
					// 设置状态检查定时器，更频繁地检查用户在线状态
					statusTicker := time.NewTicker(5 * time.Second)
					defer statusTicker.Stop()

					for {
						select {
						case <-statusTicker.C:
							if !wsManager.IsUserOnline(uid) {
								// 用户已离线，关闭停止通道
								logger.Debug("用户已离线，发送停止信号", zap.Uint("用户ID", uid))
								close(stopChan)
								return
							}
						}
					}
				}()

				// 主循环
				for {
					select {
					case <-ticker.C:
						// 定期检查是否有新消息需要处理
						// 在处理前再次检查用户状态
						if wsManager.IsUserOnline(uid) {
							if err := messageService.ProcessOfflineMessages(uid); err != nil {
								logger.Error("定期处理消息失败", zap.Uint("用户ID", uid), zap.Error(err))
							}
						}
					case <-stopChan:
						// 收到停止信号，结束处理
						logger.Debug("收到停止信号，终止消息处理", zap.Uint("用户ID", uid))
						return
					}
				}
			}(userIDUint)

			// 处理离线消息
			// 使用goroutine并设置延迟，确保连接完全建立后再处理离线消息
			go func(uid uint) {
				// 给WebSocket连接一点时间建立
				time.Sleep(500 * time.Millisecond)

				// 记录处理开始
				logger.Info("开始处理用户离线消息", zap.Uint("用户ID", uid))

				// 添加panic恢复机制，防止goroutine崩溃
				defer func() {
					if r := recover(); r != nil {
						logger.Error("消息处理协程崩溃", zap.Uint("用户ID", uid), zap.Any("panic", r))
					}
				}()

				// 首先处理离线消息
				if err := messageService.ProcessOfflineMessages(uid); err != nil {
					logger.Error("处理离线消息失败", zap.Uint("用户ID", uid), zap.Error(err))
				}

				// 创建停止通道，用于安全退出
				stopChan := make(chan struct{})

				// 设置定期检查新消息的定时器
				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()

				// 启动一个监控goroutine检查用户状态
				go func() {
					// 设置状态检查定时器，更频繁地检查用户在线状态
					statusTicker := time.NewTicker(5 * time.Second)
					defer statusTicker.Stop()

					for {
						select {
						case <-statusTicker.C:
							if !wsManager.IsUserOnline(uid) {
								// 用户已离线，关闭停止通道
								logger.Debug("用户已离线，发送停止信号", zap.Uint("用户ID", uid))
								close(stopChan)
								return
							}
						}
					}
				}()

				// 主循环
				for {
					select {
					case <-ticker.C:
						// 定期检查是否有新消息需要处理
						// 在处理前再次检查用户状态
						if wsManager.IsUserOnline(uid) {
							if err := messageService.ProcessOfflineMessages(uid); err != nil {
								logger.Error("定期处理消息失败", zap.Uint("用户ID", uid), zap.Error(err))
							}
						}
					case <-stopChan:
						// 收到停止信号，结束处理
						logger.Debug("收到停止信号，终止消息处理", zap.Uint("用户ID", uid))
						return
					}
				}
			}(userIDUint)
		})
	}
}
