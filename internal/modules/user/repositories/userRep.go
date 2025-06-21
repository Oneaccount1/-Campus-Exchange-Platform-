package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*models.User, int64, error)
	ListForAdmin(page, pageSize int, search, status, startDate, endDate string) ([]*models.User, int64, error)
	UpdateStatus(userID uint, status string) error
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u *userRepository) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Delete(id uint) error {
	return u.db.Delete(&models.User{}, id).Error
}

func (u *userRepository) List(page, pageSize int) ([]*models.User, int64, error) {
	var userList []*models.User
	var tot int64
	offset := (page - 1) * pageSize
	if err := u.db.Model(&models.User{}).Count(&tot).Error; err != nil {
		return nil, 0, err
	}
	if err := u.db.Offset(offset).Limit(pageSize).Find(&userList).Error; err != nil {
		return nil, 0, err
	}
	return userList, tot, nil
}

// UpdateStatus 更新用户状态
func (u *userRepository) UpdateStatus(userID uint, status string) error {
	return u.db.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}

// ListForAdmin 管理员用户列表高级查询
func (u *userRepository) ListForAdmin(page, pageSize int, search, status, startDate, endDate string) ([]*models.User, int64, error) {
	var userList []*models.User
	var total int64

	query := u.db.Model(&models.User{})

	// 添加搜索条件
	if search != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 添加状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
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
			// 设置结束日期为当天的最后一秒
			endTime = endTime.Add(24*time.Hour - time.Second)
			query = query.Where("created_at <= ?", endTime)
		}
	}

	// 先获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 然后获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&userList).Error; err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: bootstrap.GetDB(),
	}
}
