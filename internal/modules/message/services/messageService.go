package services

import (
	"campus/internal/messaging"
	"campus/internal/models"
	"campus/internal/modules/message/api"
	"campus/internal/modules/message/repositories"
	"campus/internal/utils/errors"
	"campus/internal/websocket"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"
)

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

	// ProcessOfflineMessages 处理离线消息
	ProcessOfflineMessages(userID uint) error

	// IsUserOnline 检查用户是否在线
	IsUserOnline(userID uint) bool

	// Close 关闭服务
	Close() error

	// GetUnreadCount 获取未读消息数量
	GetUnreadCount(userID uint) (int64, error)

	// GetLastMessage 获取与联系人的最后一条消息
	GetLastMessage(userID, contactID uint) (*api.MessageResponse, error)
}

// messageService 消息服务实现
type messageService struct {
	repo          repositories.MessageRepository // 消息仓库
	wsManager     *websocket.Manager             // WebSocket管理器
	rabbitMQURL   string                         // RabbitMQ URL
	messageQueues map[uint]*messaging.RabbitMQ   // 用户消息队列缓存
	mu            sync.RWMutex                   // 读写锁
}

func (s *messageService) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// NewMessageService 创建消息服务实例
func NewMessageService(wsManager *websocket.Manager, rabbitMQURL string) MessageService {
	return &messageService{
		repo:          repositories.NewMessageRepository(),
		wsManager:     wsManager,
		rabbitMQURL:   rabbitMQURL,
		messageQueues: make(map[uint]*messaging.RabbitMQ),
		mu:            sync.RWMutex{},
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

	// 保存消息到数据库
	if err := s.repo.Create(message); err != nil {
		return nil, errors.NewInternalServerError("消息保存失败", err)
	}

	// 转换为响应格式
	messageResponse := api.ToMessageResponse(message)

	// 将消息发送到WebSocket（如果接收者在线）
	if s.wsManager.IsUserOnline(req.ReceiverID) {
		// 尝试通过WebSocket发送
		messageJSON, err := json.Marshal(messageResponse)
		if err != nil {
			log.Printf("消息序列化失败: %v", err)
			// 序列化失败，使用RabbitMQ作为备选
			s.sendViaRabbitMQ(req.ReceiverID, message)
		} else {
			sent := s.wsManager.SendMessage(req.ReceiverID, messageJSON)
			if !sent {
				// WebSocket发送失败，使用RabbitMQ
				s.sendViaRabbitMQ(req.ReceiverID, message)
			}
		}
	} else {
		// 接收者离线，使用RabbitMQ
		s.sendViaRabbitMQ(req.ReceiverID, message)
	}

	return &messageResponse, nil
}

// sendViaRabbitMQ 通过RabbitMQ发送消息
func (s *messageService) sendViaRabbitMQ(userID uint, message *models.Message) {
	// 获取或创建用户的消息队列
	rmq, err := s.getUserMessageQueue(userID)
	if err != nil {
		log.Printf("获取用户消息队列失败: %v", err)
		return
	}

	// 发布消息到队列
	if err := rmq.Publish(message); err != nil {
		log.Printf("发布消息到RabbitMQ失败: %v", err)
	}
}

// getUserMessageQueue 获取用户消息队列
func (s *messageService) getUserMessageQueue(userID uint) (*messaging.RabbitMQ, error) {
	s.mu.RLock()
	rmq, ok := s.messageQueues[userID]
	s.mu.RUnlock()

	if ok {
		return rmq, nil
	}

	// 创建新的队列连接
	queueName := messaging.GetUserQueue(userID)
	routingKey := "user_" + strconv.FormatUint(uint64(userID), 10)

	rmq, err := messaging.NewRabbitMQ(
		s.rabbitMQURL,
		"message_exchange",
		queueName,
		routingKey,
	)

	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.messageQueues[userID] = rmq
	s.mu.Unlock()

	return rmq, nil
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
		contacts[i] = api.ContactResponse{
			ID:          user.ID,
			Username:    user.Username,
			Avatar:      user.Avatar,
			LastMessage: lastMessages[i],
			LastTime:    time.Unix(lastTimes[i], 0),
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

// ProcessOfflineMessages 处理离线消息
func (s *messageService) ProcessOfflineMessages(userID uint) error {
	// 获取用户消息队列
	rmq, err := s.getUserMessageQueue(userID)
	if err != nil {
		return errors.NewInternalServerError("获取消息队列失败", err)
	}

	// 消费消息
	deliveries, err := rmq.Consume()
	if err != nil {
		return errors.NewInternalServerError("消费消息失败", err)
	}

	// 处理离线消息（非阻塞）
	go func() {
		for delivery := range deliveries {
			// 解析消息
			var message models.Message
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				log.Printf("解析离线消息失败: %v", err)
				delivery.Reject(false) // 拒绝消息
				continue
			}

			// 转换为响应格式
			messageResponse := api.ToMessageResponse(&message)
			messageJSON, _ := json.Marshal(messageResponse)

			// 检查用户是否在线
			if s.wsManager.IsUserOnline(userID) {
				// 在线，通过WebSocket发送
				sent := s.wsManager.SendMessage(userID, messageJSON)
				if sent {
					delivery.Ack(false) // 确认消息
				} else {
					delivery.Reject(true) // 重新入队
				}
			} else {
				// 用户再次离线，保留消息
				delivery.Reject(true) // 重新入队
				return                // 退出goroutine
			}
		}
	}()

	return nil
}

// IsUserOnline 检查用户是否在线
func (s *messageService) IsUserOnline(userID uint) bool {
	return s.wsManager.IsUserOnline(userID)
}

// Close 关闭服务
func (s *messageService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 关闭所有RabbitMQ连接
	for _, rmq := range s.messageQueues {
		if err := rmq.Close(); err != nil {
			log.Printf("关闭RabbitMQ连接失败: %v", err)
		}
	}

	// 清空连接缓存
	s.messageQueues = make(map[uint]*messaging.RabbitMQ)

	return nil
}

// GetLastMessage 获取与联系人的最后一条消息
func (s *messageService) GetLastMessage(userID, contactID uint) (*api.MessageResponse, error) {
	// 获取最后一条消息
	message, err := s.repo.GetLastMessage(userID, contactID)
	if err != nil {
		return nil, errors.NewInternalServerError("获取最近消息失败", err)
	}

	// 转换为响应格式
	messageResponse := api.ToMessageResponse(message)
	return &messageResponse, nil
}
