package api

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`                     // 接收者ID
	Content    string `json:"content" binding:"required"`                         // 消息内容
	ProductID  uint   `json:"product_id"`                                         // 商品ID
	Type       string `json:"type" binding:"omitempty,oneof=text system product"` // 消息类型
}

// MarkReadRequest 标记消息已读请求
type MarkReadRequest struct {
	MessageIDs []uint `json:"message_ids"` // 消息IDs，为空则标记所有消息
}

// MessageQueryParams 消息查询参数
type MessageQueryParams struct {
	ContactID uint `form:"contact_id" binding:"required"` // 联系人ID
	Limit     int  `form:"limit,default=20"`              // 每页消息数量
	Offset    int  `form:"offset,default=0"`              // 偏移量
}

// WebSocketAuthRequest WebSocket认证请求
type WebSocketAuthRequest struct {
	Token string `json:"token" binding:"required"` // 认证token
}
