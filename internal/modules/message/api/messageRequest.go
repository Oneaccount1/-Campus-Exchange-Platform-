package api

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`                     // 接收者ID
	Content    string `json:"content" binding:"required"`                         // 消息内容
	ProductID  uint   `json:"product_id,omitempty"`                               // 商品ID（可选）
	Type       string `json:"type" binding:"omitempty,oneof=text system product"` // 消息类型
}

// MarkReadRequest 标记消息已读请求
type MarkReadRequest struct {
	MessageIDs []uint `json:"message_ids"` // 消息ID列表，为空则标记所有
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

// CreateConversationRequest 创建会话请求
type CreateConversationRequest struct {
	UserID    uint `json:"user_id" binding:"required"` // 用户ID
	ProductID uint `json:"product_id"`                 // 商品ID（可选）
}

// DeleteConversationRequest 删除会话请求
type DeleteConversationRequest struct {
	ContactID uint `json:"contact_id" binding:"required"` // 联系人ID
}

// AdminMessageListRequest 管理员获取消息列表请求
type AdminMessageListRequest struct {
	Page      uint   `json:"page" form:"page"`
	Size      uint   `json:"size" form:"size"`
	Search    string `json:"search" form:"search"`
	Type      string `json:"type" form:"type"` // 消息类型：user, system
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
}

// AdminConversationListRequest 管理员获取会话列表请求
type AdminConversationListRequest struct {
	Page   uint   `json:"page" form:"page"`
	Size   uint   `json:"size" form:"size"`
	Search string `json:"search" form:"search"`
}

// AdminMessageHistoryRequest 管理员获取会话消息历史请求
type AdminMessageHistoryRequest struct {
	User1ID   uint   `json:"user1_id" form:"user1_id" binding:"required"` // 用户1的ID
	User2ID   uint   `json:"user2_id" form:"user2_id" binding:"required"` // 用户2的ID
	Page      uint   `json:"page" form:"page"`                            // 页码
	Size      uint   `json:"size" form:"size"`                            // 每页数量
	StartDate string `json:"start_date" form:"start_date"`                // 开始日期
	EndDate   string `json:"end_date" form:"end_date"`                    // 结束日期
	Keyword   string `json:"keyword" form:"keyword"`                      // 内容关键词
	Status    string `json:"status" form:"status"`                        // 消息状态：已读/未读
}

// AdminSendSystemMessageRequest 管理员发送系统消息请求
type AdminSendSystemMessageRequest struct {
	ReceiverID uint   `json:"receiver_id"` // 接收者ID，0表示发送给所有用户
	Content    string `json:"content" binding:"required"`
	Title      string `json:"title"` // 可选标题
}
