package controllers

import (
	"campus/internal/modules/message/api"
	"campus/internal/modules/message/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// MessageController 消息控制器
type MessageController struct {
	service services.MessageService
}

// NewMessageController 创建消息控制器实例
func NewMessageController(service services.MessageService) *MessageController {
	return &MessageController{
		service: service,
	}
}

// SendMessage 发送消息
func (c *MessageController) SendMessage(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 绑定请求
	var req api.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 发送消息
	result, err := c.service.SendMessage(userID.(uint), req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "消息发送成功", result)
}

// GetMessages 获取与特定联系人的消息
func (c *MessageController) GetMessages(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 获取联系人ID
	contactIDStr := ctx.Param("contactId")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效的联系人ID", err))
		return
	}

	// 获取分页参数
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	// 请求消息列表
	result, err := c.service.GetMessagesByContact(userID.(uint), uint(contactID), limit, offset)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, result)
}

// GetContacts 获取联系人列表
func (c *MessageController) GetContacts(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 获取联系人列表
	result, err := c.service.GetContacts(userID.(uint))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, result)
}

// MarkAsRead 将消息标记为已读
func (c *MessageController) MarkAsRead(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 获取联系人ID
	contactIDStr := ctx.Param("contactId")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效的联系人ID", err))
		return
	}

	// 绑定请求
	var req api.MarkReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果没有提供消息ID列表，则标记所有消息为已读
		req.MessageIDs = []uint{}
	}

	// 标记消息为已读
	if err := c.service.MarkMessagesAsRead(userID.(uint), uint(contactID), req.MessageIDs); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "消息已标记为已读", nil)
}

// HandleWebSocket 处理WebSocket连接
func (c *MessageController) HandleWebSocket(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "未授权"})
		return
	}

	// 处理WebSocket连接
	c.service.ProcessOfflineMessages(userID.(uint))
}
