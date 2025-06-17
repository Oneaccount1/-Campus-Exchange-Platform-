package services

import (
	"campus/internal/messaging"
	"campus/internal/models"
	"campus/internal/modules/message/api"
	"campus/internal/modules/message/repositories"
	"campus/internal/websocket"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// 模拟消息仓库
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message *models.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessages(userID, contactID uint, limit, offset int) ([]models.Message, int64, error) {
	args := m.Called(userID, contactID, limit, offset)
	return args.Get(0).([]models.Message), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageRepository) MarkAsRead(messageIDs []uint, userID uint) error {
	args := m.Called(messageIDs, userID)
	return args.Error(0)
}

func (m *MockMessageRepository) MarkAllAsRead(userID, contactID uint) error {
	args := m.Called(userID, contactID)
	return args.Error(0)
}

func (m *MockMessageRepository) GetContactList(userID uint) ([]models.User, []int64, []string, []int64, []uint, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Get(1).([]int64), args.Get(2).([]string), args.Get(3).([]int64), args.Get(4).([]uint), args.Error(5)
}

func (m *MockMessageRepository) GetUnreadCount(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageRepository) Delete(messageID, userID uint) error {
	args := m.Called(messageID, userID)
	return args.Error(0)
}

func (m *MockMessageRepository) Withdraw(messageID, userID uint) error {
	args := m.Called(messageID, userID)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByID(messageID uint) (*models.Message, error) {
	args := m.Called(messageID)
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageRepository) GetLastMessage(userID, contactID uint) (*models.Message, error) {
	args := m.Called(userID, contactID)
	return args.Get(0).(*models.Message), args.Error(1)
}

// 模拟WebSocket管理器
type MockWebSocketManager struct {
	mock.Mock
}

func (m *MockWebSocketManager) IsUserOnline(userID uint) bool {
	args := m.Called(userID)
	return args.Bool(0)
}

func (m *MockWebSocketManager) SendMessage(userID uint, message []byte) bool {
	args := m.Called(userID, message)
	return args.Bool(0)
}

func (m *MockWebSocketManager) HandleConnection(writer interface{}, request interface{}, userID uint) {
	m.Called(writer, request, userID)
}

// 测试发送消息
func TestSendMessage(t *testing.T) {
	// 创建模拟对象
	mockRepo := new(MockMessageRepository)
	mockWsManager := new(MockWebSocketManager)

	// 创建服务实例
	service := &messageService{
		repo:          mockRepo,
		wsManager:     mockWsManager,
		rabbitMQURL:   "",
		messageQueues: make(map[uint]*messaging.RabbitMQ),
		mu:            sync.RWMutex{},
	}

	// 设置测试数据
	senderID := uint(1)
	receiverID := uint(2)

	// 创建请求
	req := api.SendMessageRequest{
		ReceiverID: receiverID,
		Content:    "Hello, this is a test message",
		Type:       "text",
	}

	// 预期创建的消息
	expectedMessage := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    req.Content,
		IsRead:     false,
	}

	// 设置模拟行为
	mockRepo.On("Create", mock.AnythingOfType("*models.Message")).Return(nil).Run(func(args mock.Arguments) {
		message := args.Get(0).(*models.Message)
		message.ID = 1 // 模拟数据库自动生成ID
		message.CreatedAt = time.Now()
	})

	// 设置接收者在线
	mockWsManager.On("IsUserOnline", receiverID).Return(true)
	mockWsManager.On("SendMessage", receiverID, mock.Anything).Return(true)

	// 执行测试
	messageResponse, err := service.SendMessage(senderID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, messageResponse)
	assert.Equal(t, senderID, messageResponse.SenderID)
	assert.Equal(t, receiverID, messageResponse.ReceiverID)
	assert.Equal(t, req.Content, messageResponse.Content)
	assert.False(t, messageResponse.IsRead)

	// 验证模拟对象的调用
	mockRepo.AssertExpectations(t)
	mockWsManager.AssertExpectations(t)
}

// 测试获取最后一条消息
func TestGetLastMessage(t *testing.T) {
	// 创建模拟对象
	mockRepo := new(MockMessageRepository)
	mockWsManager := new(MockWebSocketManager)

	// 创建服务实例
	service := &messageService{
		repo:          mockRepo,
		wsManager:     mockWsManager,
		rabbitMQURL:   "",
		messageQueues: make(map[uint]*messaging.RabbitMQ),
		mu:            sync.RWMutex{},
	}

	// 设置测试数据
	userID := uint(1)
	contactID := uint(2)

	// 模拟返回的消息
	lastMessage := &models.Message{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
		},
		SenderID:   contactID,
		ReceiverID: userID,
		Content:    "Last test message",
		IsRead:     true,
	}

	// 设置模拟行为
	mockRepo.On("GetLastMessage", userID, contactID).Return(lastMessage, nil)

	// 执行测试
	messageResponse, err := service.GetLastMessage(userID, contactID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, messageResponse)
	assert.Equal(t, lastMessage.ID, messageResponse.ID)
	assert.Equal(t, lastMessage.SenderID, messageResponse.SenderID)
	assert.Equal(t, lastMessage.ReceiverID, messageResponse.ReceiverID)
	assert.Equal(t, lastMessage.Content, messageResponse.Content)
	assert.Equal(t, lastMessage.IsRead, messageResponse.IsRead)

	// 验证模拟对象的调用
	mockRepo.AssertExpectations(t)
}
