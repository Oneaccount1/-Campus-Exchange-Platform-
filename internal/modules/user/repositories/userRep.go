package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*models.User, int64, error)
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

func NewUserRepository() UserRepository {
	return &userRepository{
		db: bootstrap.GetDB(),
	}
}
