package services

import (
	"campus/internal/models"
	"campus/internal/modules/message/api"
	"campus/internal/modules/message/repositories"
	"campus/internal/utils/errors"
	"encoding/json"
	"log"
	"time"
)

// RabbitMQPublisher defines the interface for publishing messages to RabbitMQ.
// This allows for loose coupling and easier testing.
type RabbitMQPublisher interface {
	Publish(body []byte, contentType string) error
}

// MessageService 消息服务接口
type MessageService interface {
	// SendMessage 发送消息
	SendMessage(senderID uint, req api.SendMessageRequest) (*api.MessageResponse, error)

	// GetMessagesByContact 获取与联系人的消息
	GetMessagesByContact(userID, contactID uint, limit, offset int) (*api.MessageListResponse, error)

	// MarkMessagesAsRead 标记消息为已读
	MarkMessagesAsRead(userID uint, contactID uint, messageIDs []uint) error

	// GetContacts 获取联系人列表
	GetContacts(userID uint) (*api.ContactListResponse, error)

	// GetUnreadCount 获取未读消息数量
	GetUnreadCount(userID uint) (int64, error)

	// GetLastMessage 获取与联系人的最后一条消息
	GetLastMessage(userID, contactID uint) (*api.MessageResponse, error)
}

// messageService 消息服务实现
type messageService struct {
	repo      repositories.MessageRepository // 消息仓库
	publisher RabbitMQPublisher              // RabbitMQ a publisher
}

func (s *messageService) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// NewMessageService 创建消息服务实例
func NewMessageService(repo repositories.MessageRepository, publisher RabbitMQPublisher) MessageService {
	return &messageService{
		repo:      repo,
		publisher: publisher,
	}
}

// SendMessage 发送消息
func (s *messageService) SendMessage(senderID uint, req api.SendMessageRequest) (*api.MessageResponse, error) {
	// 检查发送者和接收者是否相同
	if senderID == req.ReceiverID {
		return nil, errors.NewBadRequestError("不能给自己发送消息", nil)
	}

	// 创建消息
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		ProductID:  req.ProductID,
		IsRead:     false,
	}

	// 1. 保存消息到数据库
	if err := s.repo.Create(message); err != nil {
		return nil, errors.NewInternalServerError("消息保存失败", err)
	}

	// 转换为响应格式
	messageResponse := api.ToMessageResponse(message)

	// 2. 将消息发布到RabbitMQ，由后台消费者处理推送
	messageJSON, err := json.Marshal(messageResponse)
	if err != nil {
		log.Printf("消息序列化失败: %v", err)
		// Even if publishing fails, the message is in DB. Return success.
		return &messageResponse, nil
	}

	if err := s.publisher.Publish(messageJSON, "application/json"); err != nil {
		log.Printf("通过RabbitMQ发布消息失败: %v", err)
		// Even if publishing fails, the message is in DB. Return success.
	} else {
		log.Printf("消息已发布到RabbitMQ，由消费者异步处理")
	}

	return &messageResponse, nil
}

// GetMessagesByContact 获取与联系人的消息
func (s *messageService) GetMessagesByContact(userID, contactID uint, limit, offset int) (*api.MessageListResponse, error) {
	// 获取消息列表
	messages, total, err := s.repo.GetMessages(userID, contactID, limit, offset)
	if err != nil {
		return nil, errors.NewInternalServerError("获取消息失败", err)
	}

	// 自动标记为已读
	if err := s.repo.MarkAllAsRead(userID, contactID); err != nil {
		log.Printf("标记消息为已读失败: %v", err)
	}

	// 构建响应
	response := &api.MessageListResponse{
		Total:    int(total),
		Messages: api.ToMessageResponseList(messages),
	}

	return response, nil
}

// MarkMessagesAsRead 标记消息为已读
func (s *messageService) MarkMessagesAsRead(userID uint, contactID uint, messageIDs []uint) error {
	// 如果messageIDs为空，则标记所有消息已读
	if len(messageIDs) == 0 {
		return s.repo.MarkAllAsRead(userID, contactID)
	}

	// 否则标记指定消息已读
	return s.repo.MarkAsRead(messageIDs, userID)
}

// GetContacts 获取联系人列表
func (s *messageService) GetContacts(userID uint) (*api.ContactListResponse, error) {
	// 获取联系人原始数据
	users, unreadCounts, lastMessages, lastTimes, productIDs, err := s.repo.GetContactList(userID)
	if err != nil {
		return nil, errors.NewInternalServerError("获取联系人列表失败", err)
	}

	// 组装联系人列表
	contacts := make([]api.ContactResponse, len(users))
	for i, user := range users {
		// 将浮点数时间戳转换为time.Time，处理秒和毫秒部分
		seconds := int64(lastTimes[i])
		nanoseconds := int64((lastTimes[i] - float64(seconds)) * 1e9)
		lastTime := time.Unix(seconds, nanoseconds)

		contacts[i] = api.ContactResponse{
			ID:          user.ID,
			Username:    user.Username,
			Avatar:      user.Avatar,
			LastMessage: lastMessages[i],
			LastTime:    lastTime,
			Unread:      int(unreadCounts[i]),
			ProductID:   productIDs[i],
		}
	}

	// 构建响应
	response := &api.ContactListResponse{
		Total:    len(contacts),
		Contacts: contacts,
	}

	return response, nil
}

// GetLastMessage 获取与联系人的最后一条消息
func (s *messageService) GetLastMessage(userID, contactID uint) (*api.MessageResponse, error) {
	message, err := s.repo.GetLastMessage(userID, contactID)
	if err != nil {
		return nil, errors.NewNotFoundError("未找到消息", err)
	}
	response := api.ToMessageResponse(message)
	return &response, nil
}
