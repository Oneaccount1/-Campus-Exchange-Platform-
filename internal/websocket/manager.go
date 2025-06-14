package websocket

import (
	"campus/internal/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
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
			}
			m.Clients[clientReg.UserID] = clientReg.Conn
			m.ClientMux.Unlock()
			log.Printf("用户 %d 已连接WebSocket", clientReg.UserID)
		case UnregisterUserID := <-m.Unregister:
			m.ClientMux.Lock()
			if conn, ok := m.Clients[UnregisterUserID]; ok {
				close(conn.Send)
				delete(m.Clients, UnregisterUserID)
				log.Printf("用户 %d 已断开WebSocket连接", UnregisterUserID)
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
	// 这里应添加消息转JSON代码
	// ...
	return false
}

// 处理WebSocket连接
func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request, userID uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
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
				log.Printf("WebSocket错误: %v", err)
				break
			}
		}
		// 处理收到的消息(后续会添加)
		log.Printf("从用户 %d 收到WebSocket消息: %s", userID, string(message))
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
					log.Printf("通道关闭, err : %v", err)
					return
				}
			}
			// 获取 Writer
			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err := writer.Write(message); err != nil {
				log.Printf("写入消息错误 userID %d, err : %v", userID, err)
				return
			}
			if err := writer.Close(); err != nil {
				log.Printf("关闭写连接失败 err : %v", err)
				return
			}

		}
	}
}
