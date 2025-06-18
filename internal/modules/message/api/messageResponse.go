package api

import (
	"campus/internal/models"
	"time"
)

// MessageResponse 消息响应
type MessageResponse struct {
	ID         uint      `json:"id"`          // 消息ID
	SenderID   uint      `json:"sender_id"`   // 发送者ID
	ReceiverID uint      `json:"receiver_id"` // 接收者ID
	Content    string    `json:"content"`     // 内容
	IsRead     bool      `json:"is_read"`     // 是否已读
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	ProductID  uint      `json:"product_id"`  // 商品ID
}

// ContactResponse 联系人响应
type ContactResponse struct {
	ID          uint      `json:"id"`           // 用户ID
	Username    string    `json:"username"`     // 用户名
	Avatar      string    `json:"avatar"`       // 头像
	LastMessage string    `json:"last_message"` // 最后一条消息
	LastTime    time.Time `json:"last_time"`    // 最后消息时间
	Unread      int       `json:"unread"`       // 未读消息数
	ProductID   uint      `json:"product_id"`   // 商品ID
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Total    int               `json:"total"`    // 总消息数
	Messages []MessageResponse `json:"messages"` // 消息列表
}

// ContactListResponse 联系人列表响应
type ContactListResponse struct {
	Total    int               `json:"total"`    // 总联系人数
	Contacts []ContactResponse `json:"contacts"` // 联系人列表
}

// ConversationResponse 会话响应
type ConversationResponse struct {
	ID       uint   `json:"id"`       // 联系人ID
	Username string `json:"username"` // 联系人用户名
	Avatar   string `json:"avatar"`   // 联系人头像
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
		ID:          contact.UserID,
		Username:    contact.Username,
		Avatar:      contact.Avatar,
		LastMessage: contact.LastMessage,
		LastTime:    contact.LastTime,
		Unread:      contact.UnreadCount,
		ProductID:   contact.ProductID,
	}
}
