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

	// 管理员接口
	GetMessagesForAdmin(req *api.AdminMessageListRequest) (*api.AdminMessageListResponse, error)
	GetConversationsForAdmin(req *api.AdminConversationListRequest) (*api.AdminConversationListResponse, error)
	GetMessageHistoryForAdmin(req *api.AdminMessageHistoryRequest) (*api.AdminMessageHistoryResponse, error)
	SendSystemMessage(req *api.AdminSendSystemMessageRequest) error
	DeleteMessage(messageID uint) error
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
	users, unreadCounts, lastMessages, lastTimes, _, err := s.repo.GetContactList(userID)
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
			ID:           user.ID,
			Username:     user.Username,
			Avatar:       user.Avatar,
			LastMessage:  lastMessages[i],
			LastTime:     lastTime,
			UnreadCount:  int(unreadCounts[i]),
			ProductCount: 0, // 设置默认值
		}
	}

	// 构建响应
	response := &api.ContactListResponse{
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

// GetMessagesForAdmin 管理员获取消息列表
func (s *messageService) GetMessagesForAdmin(req *api.AdminMessageListRequest) (*api.AdminMessageListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	// 获取消息列表
	messages, total, err := s.repo.GetMessagesForAdmin(
		req.Search, req.Type, req.StartDate, req.EndDate, req.Page, req.Size)
	if err != nil {
		return nil, errors.NewInternalServerError("获取消息列表失败", err)
	}

	// 构建响应
	response := &api.AdminMessageListResponse{
		Total: total,
		List:  make([]api.AdminMessageItem, 0, len(messages)),
	}

	// 填充消息数据
	for _, msg := range messages {
		// 确定消息类型
		msgType := "user"
		if msg.SenderID == 0 {
			msgType = "system"
		}

		// 确定发送者和接收者名称
		var senderName, receiverName string

		if msg.SenderID == 0 {
			senderName = "系统"
		} else if msg.Sender.ID > 0 {
			senderName = msg.Sender.Username
		}

		if msg.Receiver.ID > 0 {
			receiverName = msg.Receiver.Username
		}

		// 确定消息状态
		status := "未读"
		if msg.IsRead {
			status = "已读"
		}

		// 构建消息项
		item := api.AdminMessageItem{
			ID:         msg.ID,
			Type:       msgType,
			SenderID:   msg.SenderID,
			Sender:     senderName,
			ReceiverID: msg.ReceiverID,
			Receiver:   receiverName,
			Content:    msg.Content,
			CreateTime: msg.CreatedAt,
			Status:     status,
		}

		// 如果消息已读，添加阅读时间
		if msg.IsRead && !msg.ReadTime.IsZero() {
			item.ReadTime = msg.ReadTime
		}

		response.List = append(response.List, item)
	}

	return response, nil
}

// GetConversationsForAdmin 管理员获取会话列表
func (s *messageService) GetConversationsForAdmin(req *api.AdminConversationListRequest) (*api.AdminConversationListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	// 获取会话列表
	conversations, total, err := s.repo.GetConversationsForAdmin(req.Search, req.Page, req.Size)
	if err != nil {
		return nil, errors.NewInternalServerError("获取会话列表失败", err)
	}

	// 构建响应
	response := &api.AdminConversationListResponse{
		Total: total,
		List:  make([]api.AdminConversationItem, 0, len(conversations)),
	}

	// 填充会话数据
	for _, conv := range conversations {
		// 构建会话项
		item := api.AdminConversationItem{
			ID:          conv.ID,
			User1ID:     conv.User1ID,
			User1Name:   conv.User1Name,
			User1Avatar: conv.User1Avatar,
			User2ID:     conv.User2ID,
			User2Name:   conv.User2Name,
			User2Avatar: conv.User2Avatar,
			LastMessage: conv.LastMessage,
			LastTime:    conv.LastTime,
			UnreadCount: conv.UnreadCount,
		}

		response.List = append(response.List, item)
	}

	return response, nil
}

// GetMessageHistoryForAdmin 管理员获取会话消息历史
func (s *messageService) GetMessageHistoryForAdmin(req *api.AdminMessageHistoryRequest) (*api.AdminMessageHistoryResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}

	// 获取消息历史
	messages, total, err := s.repo.GetMessageHistoryForAdmin(req.User1ID, req.User2ID, req.Page, req.Size)
	if err != nil {
		return nil, errors.NewInternalServerError("获取消息历史失败", err)
	}

	// 构建响应
	response := &api.AdminMessageHistoryResponse{
		Total: total,
		List:  make([]api.AdminMessageHistoryItem, 0, len(messages)),
	}

	// 填充消息数据
	for _, msg := range messages {
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

		// 构建消息项
		item := api.AdminMessageHistoryItem{
			ID:           msg.ID,
			SenderID:     msg.SenderID,
			Sender:       senderName,
			SenderAvatar: senderAvatar,
			ReceiverID:   msg.ReceiverID,
			Content:      msg.Content,
			CreateTime:   msg.CreatedAt,
			Status:       status,
		}

		response.List = append(response.List, item)
	}

	return response, nil
}

// SendSystemMessage 管理员发送系统消息
func (s *messageService) SendSystemMessage(req *api.AdminSendSystemMessageRequest) error {
	// 如果receiverID为0，表示发送给所有用户
	if req.ReceiverID == 0 {
		// 查询所有用户ID
		// 这里简化处理，实际应该通过用户服务获取所有用户ID
		// 由于没有直接访问用户仓库的权限，这里假设已经有所有用户ID
		// 在实际应用中，应该通过用户服务获取所有用户ID

		// 假设这里有一个获取所有用户ID的方法
		// userIDs, err := s.userService.GetAllUserIDs()
		// if err != nil {
		//     return errors.NewInternalServerError("获取用户列表失败", err)
		// }

		// 这里简化处理，直接返回成功
		// 实际应该遍历所有用户ID，为每个用户发送系统消息
		// for _, userID := range userIDs {
		//     if err := s.repo.CreateSystemMessage(userID, req.Content, req.Title); err != nil {
		//         log.Printf("向用户 %d 发送系统消息失败: %v", userID, err)
		//     }
		// }

		// 简化处理，创建一个接收者ID为0的系统消息，表示发送给所有用户
		return s.repo.CreateSystemMessage(0, req.Content, req.Title)
	}

	// 发送给特定用户
	return s.repo.CreateSystemMessage(req.ReceiverID, req.Content, req.Title)
}

// DeleteMessage 删除消息
func (s *messageService) DeleteMessage(messageID uint) error {
	// 获取消息
	_, err := s.repo.GetByID(messageID)
	if err != nil {
		return errors.NewNotFoundError("消息", err)
	}

	// 删除消息（管理员可以删除任何消息）
	if err := s.repo.Delete(messageID, 0); err != nil {
		return errors.NewInternalServerError("删除消息失败", err)
	}

	return nil
}
