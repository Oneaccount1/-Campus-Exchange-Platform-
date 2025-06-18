package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
)

// FavoriteRepository 收藏仓储接口
type FavoriteRepository interface {
	// Create 创建收藏记录
	Create(favorite *models.Favorite) error
	// Delete 删除收藏记录
	Delete(userID, productID uint) error
	// GetByUser 获取用户的所有收藏
	GetByUser(userID uint, page, size uint) ([]*models.Favorite, int64, error)
	// CheckIsFavorite 检查是否已收藏
	CheckIsFavorite(userID, productID uint) (bool, error)
}

// FavoriteRepositoryImpl 收藏仓储实现
type FavoriteRepositoryImpl struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓储实例
func NewFavoriteRepository() FavoriteRepository {
	return &FavoriteRepositoryImpl{
		db: bootstrap.GetDB(),
	}
}

// Create 创建收藏记录
func (r *FavoriteRepositoryImpl) Create(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

// Delete 删除收藏记录
func (r *FavoriteRepositoryImpl) Delete(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Favorite{}).Error
}

// GetByUser 获取用户的所有收藏
func (r *FavoriteRepositoryImpl) GetByUser(userID uint, page, size uint) ([]*models.Favorite, int64, error) {
	var favorites []*models.Favorite
	var total int64

	err := r.db.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.Where("user_id = ?", userID).
		Preload("Product").
		Preload("Product.ProductImages").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&favorites).Error

	return favorites, total, err
}

// CheckIsFavorite 检查是否已收藏
func (r *FavoriteRepositoryImpl) CheckIsFavorite(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error

	return count > 0, err
}
