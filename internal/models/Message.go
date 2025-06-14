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

// Contact 联系人模型（视图模型，不直接映射到数据库）
type Contact struct {
	UserID       uint      `json:"id"`             // 用户ID
	Username     string    `json:"username"`       // 用户名
	Avatar       string    `json:"avatar"`         // 头像
	LastMessage  string    `json:"last_message"`   // 最后一条消息
	LastTime     time.Time `json:"last_time"`      // 最后消息时间
	UnreadCount  int       `json:"unread"`         // 未读消息数量
	ProductID    uint      `json:"product_id"`     // 相关商品ID
	LastSenderID uint      `json:"last_sender_id"` // 最后发送消息的用户ID
}
