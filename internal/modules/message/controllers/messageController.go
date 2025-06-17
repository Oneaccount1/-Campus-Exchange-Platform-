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

// GetUnreadCount 获取未读消息数量
func (c *MessageController) GetUnreadCount(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 获取未读消息数量
	count, err := c.service.GetUnreadCount(userID.(uint))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	// 返回未读消息数量
	response.Success(ctx, gin.H{
		"count": count,
	})
}

// CreateConversation 创建新会话
func (c *MessageController) CreateConversation(ctx *gin.Context) {
	// 验证用户是否已登录
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// 绑定请求
	var req api.CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 验证请求
	if req.UserID == 0 {
		response.HandleError(ctx, errors.NewBadRequestError("联系人ID不能为空", nil))
		return
	}

	// 不能和自己创建会话
	if currentUserID.(uint) == req.UserID {
		response.HandleError(ctx, errors.NewBadRequestError("不能与自己创建会话", nil))
		return
	}

	// 查询用户信息 - 这里应该调用用户服务获取用户信息
	// 在实际应用中，应该通过用户服务获取真实用户信息
	// 这里简化处理，直接返回基本信息
	conversationResponse := api.ConversationResponse{
		ID:       req.UserID,
		Username: "用户" + strconv.FormatUint(uint64(req.UserID), 10),
		Avatar:   "/static/default_avatar.png", // 使用默认头像
	}

	response.Success(ctx, conversationResponse)
}

// GetLastMessage 获取最后一条消息
func (c *MessageController) GetLastMessage(ctx *gin.Context) {
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

	// 获取最后一条消息
	message, err := c.service.GetLastMessage(userID.(uint), uint(contactID))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, message)
}
