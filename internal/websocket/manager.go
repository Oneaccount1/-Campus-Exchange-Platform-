package websocket

import (
	"campus/internal/models"
	"campus/internal/utils/logger"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有cros

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Connection 表示websocket连接的包装
type Connection struct {
	Conn *websocket.Conn
	Send chan []byte
}

// Manager 管理Websocket连接
type Manager struct {
	// 客户端连接映射表： 用户ID -> 连接
	Clients   map[uint]*Connection
	ClientMux sync.RWMutex

	// 注册和注销通道
	Register   chan *ClientRegistration
	Unregister chan uint
}

type ClientRegistration struct {
	UserID uint
	Conn   *Connection
}

func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[uint]*Connection),
		Register:   make(chan *ClientRegistration),
		Unregister: make(chan uint),
		ClientMux:  sync.RWMutex{},
	}
}

func (m *Manager) Start() {
	for {
		select {
		case clientReg := <-m.Register:
			m.ClientMux.Lock()
			// 如果该用户已经有连接 ，先关闭连接
			if conn, ok := m.Clients[clientReg.UserID]; ok {
				close(conn.Send)
				delete(m.Clients, clientReg.UserID)
				logger.Debugf("用户 %d 的旧连接已关闭", clientReg.UserID)
			}
			m.Clients[clientReg.UserID] = clientReg.Conn
			m.ClientMux.Unlock()
			logger.Info("WebSocket连接建立", zap.Uint("用户ID", clientReg.UserID))

		case userID := <-m.Unregister:
			m.ClientMux.Lock()
			if conn, ok := m.Clients[userID]; ok {
				close(conn.Send)
				delete(m.Clients, userID)
				logger.Info("WebSocket连接断开", zap.Uint("用户ID", userID))
			}
			m.ClientMux.Unlock()
		}
	}
}

// IsUserOnline 检查用户是否在线
func (m *Manager) IsUserOnline(userID uint) bool {
	m.ClientMux.RLock()
	_, exists := m.Clients[userID]
	m.ClientMux.RUnlock()
	return exists
}

// SendMessage 向指定用户发送消息
func (m *Manager) SendMessage(userID uint, message []byte) bool {
	m.ClientMux.RLock()
	conn, exists := m.Clients[userID]
	m.ClientMux.RUnlock()
	if exists {
		conn.Send <- message
		return true
	}
	return false
}

// SendMessageToUser 发送消息模型到指定用户
func (m *Manager) SendMessageToUser(message *models.Message) bool {
	// 将消息转换为JSON格式
	messageResponse := struct {
		ID         uint      `json:"id"`
		SenderID   uint      `json:"sender_id"`
		ReceiverID uint      `json:"receiver_id"`
		Content    string    `json:"content"`
		IsRead     bool      `json:"is_read"`
		CreatedAt  time.Time `json:"created_at"`
		ProductID  uint      `json:"product_id"`
	}{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		IsRead:     message.IsRead,
		CreatedAt:  message.CreatedAt,
		ProductID:  message.ProductID,
	}

	messageJSON, err := json.Marshal(messageResponse)
	if err != nil {
		logger.Error("消息序列化失败",
			zap.Uint("接收者ID", message.ReceiverID),
			zap.Error(err))
		return false
	}

	// 发送消息到接收者
	return m.SendMessage(message.ReceiverID, messageJSON)
}

// HandleConnection 处理WebSocket连接
func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request, userID uint) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WebSocket连接升级失败",
			zap.Uint("用户ID", userID),
			zap.Error(err))
		return
	}

	// 连接成功， 创建用户连接
	client := &Connection{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// 注册连接
	m.Register <- &ClientRegistration{
		UserID: userID,
		Conn:   client,
	}

	// 启动消息读取与写入协程
	go m.readPump(client, userID)
	go m.writePump(client, userID)

	logger.Debug("WebSocket处理协程已启动", zap.Uint("用户ID", userID))
}

func (m *Manager) readPump(c *Connection, userID uint) {
	defer func() {
		m.Unregister <- userID
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn("WebSocket意外关闭",
					zap.Uint("用户ID", userID),
					zap.Error(err))
			}
			break
		}

		// 处理ping消息
		if string(message) == "ping" {
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				logger.Warn("发送pong响应失败",
					zap.Uint("用户ID", userID),
					zap.Error(err))
			}
			continue
		}

		// 只在调试级别记录收到的消息
		logger.Debug("收到WebSocket消息",
			zap.Uint("用户ID", userID),
			zap.String("消息", string(message)))
	}
}

func (m *Manager) writePump(c *Connection, userID uint) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// 通道关闭
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Debug("关闭WebSocket连接失败",
						zap.Uint("用户ID", userID),
						zap.Error(err))
				}
				return
			}
			// 获取 Writer
			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Warn("获取WebSocket Writer失败",
					zap.Uint("用户ID", userID),
					zap.Error(err))
				return
			}

			if _, err := writer.Write(message); err != nil {
				logger.Warn("WebSocket消息写入失败",
					zap.Uint("用户ID", userID),
					zap.Error(err))
				return
			}

			if err := writer.Close(); err != nil {
				logger.Debug("关闭WebSocket Writer失败",
					zap.Uint("用户ID", userID),
					zap.Error(err))
				return
			}
		}
	}
}
