package models

import (
	"gorm.io/gorm"
	"time"
)

// 消息类型枚举
const (
	MessageTypeText    = "text"    // 文本消息
	MessageTypeSystem  = "system"  // 系统消息
	MessageTypeProduct = "product" // 商品相关消息
)

// Message 消息模型
type Message struct {
	gorm.Model
	SenderID    uint      `gorm:"not null;index" json:"sender_id"`           // 发送者ID
	Sender      User      `gorm:"foreignKey:SenderID" json:"sender"`         // 发送者
	ReceiverID  uint      `gorm:"not null;index" json:"receiver_id"`         // 接收者ID
	Receiver    User      `gorm:"foreignKey:ReceiverID" json:"receiver"`     // 接收者
	Content     string    `gorm:"size:1000;not null" json:"content"`         // 消息内容
	Type        string    `gorm:"size:20;not null;default:text" json:"type"` // 消息类型
	IsRead      bool      `gorm:"default:false" json:"is_read"`              // 是否已读
	ReadTime    time.Time `gorm:"default:null" json:"read_time"`             // 阅读时间
	ProductID   uint      `gorm:"index" json:"product_id"`                   // 相关商品ID
	Product     Product   `gorm:"foreignKey:ProductID" json:"product"`       // 相关商品
	IsDeleted   bool      `gorm:"default:false" json:"is_deleted"`           // 软删除标记
	IsWithdrawn bool      `gorm:"default:false" json:"is_withdrawn"`         // 是否已撤回
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// Contact 联系人模型（非数据库表，用于API响应）
type Contact struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Avatar       string    `json:"avatar"`
	LastMessage  string    `json:"last_message"`
	LastTime     time.Time `json:"last_time"`
	UnreadCount  int       `json:"unread_count"`
	ProductCount int       `json:"product_count"`
}

// Conversation 会话模型（用于管理员API）
type Conversation struct {
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

// MessageLog 系统消息日志
type MessageLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"type:text" json:"content"`
	Title      string    `gorm:"size:100" json:"title"`
	IsSystem   bool      `gorm:"default:false" json:"is_system"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
