package api

import (
	"campus/internal/models"
	"time"
)

// MessageResponse 消息响应
type MessageResponse struct {
	ID         uint      `json:"id"`                   // 消息ID
	SenderID   uint      `json:"sender_id"`            // 发送者ID
	ReceiverID uint      `json:"receiver_id"`          // 接收者ID
	Content    string    `json:"content"`              // 内容
	IsRead     bool      `json:"is_read"`              // 是否已读
	CreatedAt  time.Time `json:"created_at"`           // 创建时间
	ProductID  uint      `json:"product_id,omitempty"` // 商品ID
}

// ContactResponse 联系人响应
type ContactResponse struct {
	ID           uint      `json:"id"`                      // 用户ID
	Username     string    `json:"username"`                // 用户名
	Avatar       string    `json:"avatar"`                  // 头像
	LastMessage  string    `json:"last_message"`            // 最后一条消息
	LastTime     time.Time `json:"last_time"`               // 最后消息时间
	UnreadCount  int       `json:"unread_count"`            // 未读消息数
	ProductCount int       `json:"product_count,omitempty"` // 商品数量
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Total    int               `json:"total"`    // 总消息数
	Messages []MessageResponse `json:"messages"` // 消息列表
}

// ContactListResponse 联系人列表响应
type ContactListResponse struct {
	Contacts []ContactResponse `json:"contacts"` // 联系人列表
}

// ConversationResponse 会话响应
type ConversationResponse struct {
	ID       uint   `json:"id"`       // 会话ID
	Username string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 头像
}

// MessageHistoryResponse 消息历史响应
type MessageHistoryResponse struct {
	ID             uint      `json:"id"`                   // 消息ID
	SenderID       uint      `json:"sender_id"`            // 发送者ID
	Sender         string    `json:"sender"`               // 发送者名称
	SenderAvatar   string    `json:"sender_avatar"`        // 发送者头像
	ReceiverID     uint      `json:"receiver_id"`          // 接收者ID
	Receiver       string    `json:"receiver"`             // 接收者名称
	ReceiverAvatar string    `json:"receiver_avatar"`      // 接收者头像
	Content        string    `json:"content"`              // 消息内容
	CreateTime     time.Time `json:"create_time"`          // 创建时间
	Status         string    `json:"status"`               // 状态
	ProductID      uint      `json:"product_id,omitempty"` // 商品ID
}

// MessageHistoryListResponse 消息历史列表响应
type MessageHistoryListResponse struct {
	Total int64                    `json:"total"` // 总数
	List  []MessageHistoryResponse `json:"list"`  // 消息列表
}

// AdminMessageResponse 管理员消息响应
type AdminMessageResponse struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`        // 消息类型：user 或 system
	SenderID   uint      `json:"sender_id"`   // 发送者ID
	Sender     string    `json:"sender"`      // 发送者名称
	ReceiverID uint      `json:"receiver_id"` // 接收者ID
	Receiver   string    `json:"receiver"`    // 接收者名称
	Content    string    `json:"content"`     // 消息内容
	CreateTime time.Time `json:"create_time"` // 创建时间
	Status     string    `json:"status"`      // 状态：已读或未读
	ReadTime   time.Time `json:"read_time"`   // 阅读时间
}

// AdminMessageItem 管理员消息列表项
type AdminMessageItem struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`
	SenderID   uint      `json:"sender_id"`
	Sender     string    `json:"sender"`
	ReceiverID uint      `json:"receiver_id"`
	Receiver   string    `json:"receiver"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
	Status     string    `json:"status"`
	ReadTime   time.Time `json:"read_time,omitempty"`
}

// AdminMessageListResponse 管理员消息列表响应
type AdminMessageListResponse struct {
	Total int64              `json:"total"`
	List  []AdminMessageItem `json:"list"`
}

// AdminConversationItem 管理员会话列表项
type AdminConversationItem struct {
	ID          uint      `json:"id"`
	User1ID     uint      `json:"user1_id"`
	User1Name   string    `json:"user1_name"`
	User1Avatar string    `json:"user1_avatar"`
	User2ID     uint      `json:"user2_id"`
	User2Name   string    `json:"user2_name"`
	User2Avatar string    `json:"user2_avatar"`
	LastMessage string    `json:"last_message"`
	LastTime    time.Time `json:"last_time"`
	UnreadCount int64     `json:"unread_count"`
}

// AdminConversationListResponse 管理员会话列表响应
type AdminConversationListResponse struct {
	Total int64                   `json:"total"`
	List  []AdminConversationItem `json:"list"`
}

// AdminMessageHistoryItem 管理员消息历史项
type AdminMessageHistoryItem struct {
	ID             uint      `json:"id"`                   // 消息ID
	SenderID       uint      `json:"sender_id"`            // 发送者ID
	Sender         string    `json:"sender"`               // 发送者名称
	SenderAvatar   string    `json:"sender_avatar"`        // 发送者头像
	ReceiverID     uint      `json:"receiver_id"`          // 接收者ID
	Receiver       string    `json:"receiver"`             // 接收者名称
	ReceiverAvatar string    `json:"receiver_avatar"`      // 接收者头像
	Content        string    `json:"content"`              // 消息内容
	CreateTime     time.Time `json:"create_time"`          // 创建时间
	Status         string    `json:"status"`               // 状态
	ProductID      uint      `json:"product_id,omitempty"` // 商品ID
}

// AdminMessageHistoryResponse 管理员消息历史响应
type AdminMessageHistoryResponse struct {
	Total int64                     `json:"total"`
	List  []AdminMessageHistoryItem `json:"list"`
}

// ToMessageResponse 将Message模型转换为响应
func ToMessageResponse(msg *models.Message) MessageResponse {
	return MessageResponse{
		ID:         msg.ID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		IsRead:     msg.IsRead,
		CreatedAt:  msg.CreatedAt,
		ProductID:  msg.ProductID,
	}
}

// ToMessageResponseList 将消息模型列表转换为响应
func ToMessageResponseList(messages []models.Message) []MessageResponse {
	result := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		result[i] = ToMessageResponse(&msg)
	}
	return result
}

// ToContactResponse 将Contact模型转换为响应
func ToContactResponse(contact *models.Contact) ContactResponse {
	return ContactResponse{
		ID:           contact.UserID,
		Username:     contact.Username,
		Avatar:       contact.Avatar,
		LastMessage:  contact.LastMessage,
		LastTime:     contact.LastTime,
		UnreadCount:  contact.UnreadCount,
		ProductCount: contact.ProductCount,
	}
}

// ToMessageHistoryResponse 将Message模型转换为MessageHistoryResponse
func ToMessageHistoryResponse(msg *models.Message) MessageHistoryResponse {
	// 确定消息状态
	status := "未读"
	if msg.IsRead {
		status = "已读"
	}

	// 获取发送者信息
	var senderName, senderAvatar string
	if msg.SenderID == 0 {
		senderName = "系统"
		senderAvatar = "/static/system_avatar.png"
	} else if msg.Sender.ID > 0 {
		senderName = msg.Sender.Username
		senderAvatar = msg.Sender.Avatar
	}

	// 获取接收者信息
	var receiverName, receiverAvatar string
	if msg.ReceiverID > 0 && msg.Receiver.ID > 0 {
		receiverName = msg.Receiver.Username
		receiverAvatar = msg.Receiver.Avatar
	}

	return MessageHistoryResponse{
		ID:             msg.ID,
		SenderID:       msg.SenderID,
		Sender:         senderName,
		SenderAvatar:   senderAvatar,
		ReceiverID:     msg.ReceiverID,
		Receiver:       receiverName,
		ReceiverAvatar: receiverAvatar,
		Content:        msg.Content,
		CreateTime:     msg.CreatedAt,
		Status:         status,
		ProductID:      msg.ProductID,
	}
}

// ToMessageHistoryResponseList 将Message模型列表转换为MessageHistoryResponse列表
func ToMessageHistoryResponseList(messages []models.Message) []MessageHistoryResponse {
	result := make([]MessageHistoryResponse, len(messages))
	for i, msg := range messages {
		result[i] = ToMessageHistoryResponse(&msg)
	}
	return result
}
