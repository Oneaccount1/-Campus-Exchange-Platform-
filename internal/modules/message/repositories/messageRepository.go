package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
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
	GetContactList(userID uint) ([]models.User, []int64, []string, []int64, []uint, error)

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
}

// messageRepository 消息仓库实现
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓库实例
func NewMessageRepository() MessageRepository {
	return &messageRepository{
		db: bootstrap.GetDB(),
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

	// 查询消息总数
	err := r.db.Model(&models.Message{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, contactID, contactID, userID).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// 查询消息列表
	err = r.db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, contactID, contactID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Sender").
		Preload("Receiver").
		Preload("Product").
		Find(&messages).Error

	return messages, total, err
}

// MarkAsRead 标记特定消息为已读
func (r *messageRepository) MarkAsRead(messageIDs []uint, userID uint) error {
	if len(messageIDs) == 0 {
		return nil
	}
	return r.db.Model(&models.Message{}).
		Where("id IN ? AND receiver_id = ? AND is_read = ?", messageIDs, userID, false).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记用户与联系人间的所有消息为已读
func (r *messageRepository) MarkAllAsRead(userID, contactID uint) error {
	return r.db.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", contactID, userID, false).
		Update("is_read", true).Error
}

// GetContactList 获取联系人列表（返回原始数据，由服务层组装）
func (r *messageRepository) GetContactList(userID uint) ([]models.User, []int64, []string, []int64, []uint, error) {
	// 查找有过通信的所有用户
	var contactIds []uint

	// 查询与用户有过通信的用户ID
	if err := r.db.Model(&models.Message{}).
		Select("DISTINCT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS contact_id", userID).
		Where("(sender_id = ? OR receiver_id = ?)", userID, userID).
		Where("sender_id != receiver_id"). // 排除自己发给自己的消息
		Pluck("contact_id", &contactIds).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 如果没有联系人，返回空结果
	if len(contactIds) == 0 {
		return []models.User{}, []int64{}, []string{}, []int64{}, []uint{}, nil
	}

	// 获取所有联系人的用户信息
	var users []models.User
	if err := r.db.Where("id IN ?", contactIds).Find(&users).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 为每个联系人查询最后一条消息
	var lastMessages []string
	var lastTimes []int64
	var unreadCounts []int64
	var productIDs []uint

	for _, contactID := range contactIds {
		// 最后一条消息
		var lastMsg models.Message
		if err := r.db.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, contactID, contactID, userID,
		).Order("created_at DESC").First(&lastMsg).Error; err == nil {
			lastMessages = append(lastMessages, lastMsg.Content)
			lastTimes = append(lastTimes, lastMsg.CreatedAt.Unix())
			productIDs = append(productIDs, lastMsg.ProductID)
		} else {
			lastMessages = append(lastMessages, "")
			lastTimes = append(lastTimes, 0)
			productIDs = append(productIDs, 0)
		}

		// 未读消息数
		var unread int64
		r.db.Model(&models.Message{}).
			Where("sender_id = ? AND receiver_id = ? AND is_read = ?",
				contactID, userID, false).
			Count(&unread)
		unreadCounts = append(unreadCounts, unread)
	}

	return users, unreadCounts, lastMessages, lastTimes, productIDs, nil
}

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
