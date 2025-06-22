package repositories

import (
	"campus/internal/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// MessageRepository 消息仓库接口
type MessageRepository interface {
	// Create 创建消息
	Create(message *models.Message) error

	// GetMessages 获取消息列表
	GetMessages(userID, contactID uint, limit, offset int) ([]models.Message, int64, error)

	// MarkAsRead 标记特定消息为已读
	MarkAsRead(messageIDs []uint, userID uint) error

	// MarkAllAsRead 标记用户与联系人间的所有消息为已读
	MarkAllAsRead(userID, contactID uint) error

	// GetContactList 获取联系人列表
	GetContactList(userID uint) ([]models.User, []int64, []string, []float64, []uint, error)

	// GetUnreadCount 获取未读消息数
	GetUnreadCount(userID uint) (int64, error)

	// Delete 删除消息（软删除）
	Delete(messageID, userID uint) error

	// Withdraw 撤回消息
	Withdraw(messageID, userID uint) error

	// GetByID 获取单个消息
	GetByID(messageID uint) (*models.Message, error)

	// GetLastMessage 获取最后一条消息
	GetLastMessage(userID, contactID uint) (*models.Message, error)

	// 管理员接口
	GetMessagesForAdmin(search, msgType, startDate, endDate string, page, pageSize uint) ([]models.Message, int64, error)
	GetConversationsForAdmin(search string, page, pageSize uint) ([]models.Conversation, int64, error)
	GetMessageHistoryForAdmin(user1ID, user2ID uint, page, pageSize uint) ([]models.Message, int64, error)
	CreateSystemMessage(receiverID uint, content, title string) error
}

// messageRepository 消息仓库实现
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓库实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

// Create 创建消息
func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetMessages 获取消息列表
func (r *messageRepository) GetMessages(userID, contactID uint, limit, offset int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	// 查询条件：用户和联系人之间的消息
	condition := r.db.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, contactID, contactID, userID,
	)

	// 计算总记录数
	if err := condition.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询消息列表
	if err := condition.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// MarkAsRead 标记特定消息为已读
func (r *messageRepository) MarkAsRead(messageIDs []uint, userID uint) error {
	return r.db.Model(&models.Message{}).
		Where("id IN ? AND receiver_id = ? AND is_read = ?", messageIDs, userID, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": time.Now(),
		}).Error
}

// MarkAllAsRead 标记用户与联系人间的所有消息为已读
func (r *messageRepository) MarkAllAsRead(userID, contactID uint) error {
	return r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = ?", userID, contactID, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": time.Now(),
		}).Error
}

func (r *messageRepository) GetContactList(userID uint) ([]models.User, []int64, []string, []float64, []uint, error) {
	// 获取联系人ID列表（按最近消息时间排序）
	type ContactResult struct {
		ContactID uint
		LastTime  float64
	}
	var contactResults []ContactResult

	err := r.db.Model(&models.Message{}).
		Select(`
            CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS contact_id,
            MAX(UNIX_TIMESTAMP(created_at)) as last_time`,
			userID,
		).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("contact_id").
		Having("contact_id != ?", userID).
		Order("last_time DESC").
		Scan(&contactResults).Error

	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if len(contactResults) == 0 {
		return []models.User{}, []int64{}, []string{}, []float64{}, []uint{}, nil
	}

	// 提取ID并保持顺序
	contactIDs := make([]uint, len(contactResults))
	lastTimes := make([]float64, len(contactResults))
	for i, cr := range contactResults {
		contactIDs[i] = cr.ContactID
		lastTimes[i] = cr.LastTime
	}

	// 获取联系人详细信息（保持原始顺序）
	var users []models.User
	if err := r.db.Where("id IN ?", contactIDs).Find(&users).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 构建用户ID到索引的映射
	userIndex := make(map[uint]int)
	for i, user := range users {
		userIndex[user.ID] = i
	}

	// 初始化结果集
	userCount := len(users)
	unreadCounts := make([]int64, userCount)
	lastMessages := make([]string, userCount)
	productIDs := make([]uint, userCount)

	// 重置lastTimes（将用实际值替换）
	lastTimes = make([]float64, userCount)

	// 1. 批量查询未读消息数
	var unreads []struct {
		SenderID uint
		Count    int64
	}
	r.db.Model(&models.Message{}).
		Select("sender_id, count(*) as count").
		Where("receiver_id = ? AND is_read = ? AND sender_id IN ?", userID, false, contactIDs).
		Group("sender_id").
		Scan(&unreads)

	for _, u := range unreads {
		if idx, exists := userIndex[u.SenderID]; exists {
			unreadCounts[idx] = u.Count
		}
	}

	// 2. 批量查询最后一条消息
	var latestMessages []models.Message
	r.db.Raw(`
        SELECT m1.* 
        FROM messages m1
        INNER JOIN (
            SELECT 
                CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS contact_id,
                MAX(created_at) AS max_time
            FROM messages
            WHERE (sender_id = ? OR receiver_id = ?)
                AND (sender_id IN ? OR receiver_id IN ?)
            GROUP BY contact_id
        ) m2 ON (m1.sender_id = m2.contact_id OR m1.receiver_id = m2.contact_id)
            AND m1.created_at = m2.max_time
        `,
		userID, userID, userID, contactIDs, contactIDs,
	).Scan(&latestMessages)

	// 填充最后一条消息数据
	for _, msg := range latestMessages {
		var contactID uint
		if msg.SenderID == userID {
			contactID = msg.ReceiverID
		} else {
			contactID = msg.SenderID
		}

		if idx, exists := userIndex[contactID]; exists {
			lastMessages[idx] = msg.Content
			lastTimes[idx] = float64(msg.CreatedAt.Unix())
			productIDs[idx] = msg.ProductID
		}
	}

	return users, unreadCounts, lastMessages, lastTimes, productIDs, nil
}

//// GetContactList 获取联系人列表
//func (r *messageRepository) GetContactList(userID uint) ([]models.User, []int64, []string, []float64, []uint, error) {
//	// 查询与当前用户有过消息往来的所有用户ID
//	var contactIDs []uint
//	err := r.db.Model(&models.Message{}).
//		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS contact_id").
//		Where("sender_id = ? OR receiver_id = ?", userID, userID).
//		Group("contact_id").
//		Having("contact_id != ?", userID).
//		Pluck("contact_id", &contactIDs).Error
//
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//
//	if len(contactIDs) == 0 {
//		return []models.User{}, []int64{}, []string{}, []float64{}, []uint{}, nil
//	}
//
//	// 查询联系人信息
//	var users []models.User
//	if err := r.db.Where("id IN ?", contactIDs).Find(&users).Error; err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//
//	// 查询每个联系人的未读消息数
//	unreadCounts := make([]int64, len(contactIDs))
//	lastMessages := make([]string, len(contactIDs))
//	lastTimes := make([]float64, len(contactIDs))
//	productIDs := make([]uint, len(contactIDs))
//
//	for i, contactID := range contactIDs {
//		// 查询未读消息数
//		var unreadCount int64
//		r.db.Model(&models.Message{}).
//			Where("sender_id = ? AND receiver_id = ? AND is_read = ?", contactID, userID, false).
//			Count(&unreadCount)
//		unreadCounts[i] = unreadCount
//
//		// 查询最后一条消息
//		var lastMessage models.Message
//		if err := r.db.Where(
//			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
//			userID, contactID, contactID, userID,
//		).Order("created_at DESC").First(&lastMessage).Error; err == nil {
//			lastMessages[i] = lastMessage.Content
//			lastTimes[i] = float64(lastMessage.CreatedAt.Unix())
//			productIDs[i] = lastMessage.ProductID
//		}
//	}
//
//	return users, unreadCounts, lastMessages, lastTimes, productIDs, nil
//}

// GetUnreadCount 获取未读消息数
func (r *messageRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// Delete 删除消息（软删除）
func (r *messageRepository) Delete(messageID, userID uint) error {
	return r.db.Model(&models.Message{}).
		Where("id = ? AND (sender_id = ? OR receiver_id = ?)", messageID, userID, userID).
		Update("is_deleted", true).Error
}

// Withdraw 撤回消息
func (r *messageRepository) Withdraw(messageID, userID uint) error {
	return r.db.Model(&models.Message{}).
		Where("id = ? AND sender_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL 2 MINUTE)",
			messageID, userID).
		Update("is_withdrawn", true).Error
}

// GetByID 获取单个消息
func (r *messageRepository) GetByID(messageID uint) (*models.Message, error) {
	var message models.Message
	err := r.db.First(&message, messageID).Error
	return &message, err
}

// GetLastMessage 获取最后一条消息
func (r *messageRepository) GetLastMessage(userID, contactID uint) (*models.Message, error) {
	var message models.Message

	// 查询用户和联系人之间的最后一条消息
	// 这里的查询条件确保了只获取用户和联系人之间的消息，不管是谁发给谁的
	err := r.db.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, contactID, contactID, userID,
	).Order("created_at DESC").First(&message).Error

	return &message, err
}

// GetMessagesForAdmin 管理员获取消息列表
func (r *messageRepository) GetMessagesForAdmin(search, msgType, startDate, endDate string, page, pageSize uint) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	query := r.db.Model(&models.Message{})

	// 添加搜索条件
	if search != "" {
		query = query.Where("content LIKE ?", "%"+search+"%")
	}

	// 添加消息类型筛选
	if msgType == "user" {
		query = query.Where("sender_id > 0")
	} else if msgType == "system" {
		query = query.Where("sender_id = 0")
	}

	// 添加日期筛选
	if startDate != "" {
		startTime, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}

	if endDate != "" {
		endTime, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// 增加一天，使得结束日期是包含当天的
			endTime = endTime.AddDate(0, 0, 1)
			query = query.Where("created_at < ?", endTime)
		}
	}

	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetConversationsForAdmin 管理员获取会话列表
func (r *messageRepository) GetConversationsForAdmin(search string, page, pageSize uint) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64

	// 查询所有不同的用户对（去重）
	query := `
		WITH user_pairs AS (
			SELECT 
				CASE WHEN sender_id < receiver_id THEN sender_id ELSE receiver_id END AS user1_id,
				CASE WHEN sender_id < receiver_id THEN receiver_id ELSE sender_id END AS user2_id
			FROM messages
			WHERE sender_id > 0 AND receiver_id > 0
			GROUP BY user1_id, user2_id
		)
		SELECT 
			p.user1_id, 
			p.user2_id,
			u1.username AS user1_name,
			u1.avatar AS user1_avatar,
			u2.username AS user2_name,
			u2.avatar AS user2_avatar
		FROM user_pairs p
		JOIN users u1 ON p.user1_id = u1.id
		JOIN users u2 ON p.user2_id = u2.id
	`

	// 添加搜索条件
	if search != "" {
		query += fmt.Sprintf(" WHERE u1.username LIKE '%%%s%%' OR u2.username LIKE '%%%s%%'", search, search)
	}

	// 计算总记录数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", query)
	if err := r.db.Raw(countQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加分页
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	// 执行查询
	rows, err := r.db.Raw(query).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 处理结果
	for rows.Next() {
		var conv models.Conversation
		if err := r.db.ScanRows(rows, &conv); err != nil {
			return nil, 0, err
		}

		// 获取最后一条消息
		var lastMessage models.Message
		if err := r.db.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			conv.User1ID, conv.User2ID, conv.User2ID, conv.User1ID,
		).Order("created_at DESC").First(&lastMessage).Error; err == nil {
			conv.LastMessage = lastMessage.Content
			conv.LastTime = lastMessage.CreatedAt
		}

		// 获取未读消息数
		var unreadCount int64
		if err := r.db.Model(&models.Message{}).
			Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND is_read = ?",
				conv.User1ID, conv.User2ID, conv.User2ID, conv.User1ID, false).
			Count(&unreadCount).Error; err == nil {
			conv.UnreadCount = unreadCount
		}

		conversations = append(conversations, conv)
	}

	return conversations, total, nil
}

// GetMessageHistoryForAdmin 管理员获取会话消息历史
func (r *messageRepository) GetMessageHistoryForAdmin(user1ID, user2ID uint, page, pageSize uint) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	// 查询条件：两个用户之间的消息
	query := r.db.Model(&models.Message{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	)

	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// CreateSystemMessage 创建系统消息
func (r *messageRepository) CreateSystemMessage(receiverID uint, content, title string) error {
	// 构建系统消息
	message := models.Message{
		SenderID:   0, // 系统消息的发送者ID为0
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}

	// 如果有标题，添加到内容前面
	if title != "" {
		message.Content = fmt.Sprintf("[%s] %s", title, content)
	}

	return r.db.Create(&message).Error
}

// ToContactResponse 将查询结果转换为Contact模型
func (r *messageRepository) ToContactResponse(userID uint, username, avatar, lastMessage string, lastTime time.Time, unreadCount int) *models.Contact {
	return &models.Contact{
		UserID:       userID,
		Username:     username,
		Avatar:       avatar,
		LastMessage:  lastMessage,
		LastTime:     lastTime,
		UnreadCount:  unreadCount,
		ProductCount: 0, // 默认值，如果需要可以通过其他查询填充
	}
}
