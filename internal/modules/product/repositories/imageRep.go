package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
)

// ProductImageRepository 定义商品图片仓储接口
type ProductImageRepository interface {
	// CreateImages 批量创建商品图片记录
	CreateImages(images []*models.ProductImage) error
	// DeleteImagesByProductID 根据商品 ID 删除对应的图片记录
	DeleteImagesByProductID(productID string) error
	// GetImagesByProductID 根据商品 ID 获取对应的图片记录
	GetImagesByProductID(productID string) ([]*models.ProductImage, error)
}

// ProductImageRepositoryImpl 实现商品图片仓储接口
type ProductImageRepositoryImpl struct {
	db *gorm.DB
}

// NewProductImageRepository 创建商品图片仓储实例
func NewProductImageRepository() ProductImageRepository {
	return &ProductImageRepositoryImpl{
		db: bootstrap.GetDB(),
	}
}

// CreateImages 批量创建商品图片记录
func (r *ProductImageRepositoryImpl) CreateImages(images []*models.ProductImage) error {
	return r.db.Create(images).Error
}

// DeleteImagesByProductID 根据商品 ID 删除对应的图片记录
func (r *ProductImageRepositoryImpl) DeleteImagesByProductID(productID string) error {
	return r.db.Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error
}

// GetImagesByProductID 根据商品 ID 获取对应的图片记录
func (r *ProductImageRepositoryImpl) GetImagesByProductID(productID string) ([]*models.ProductImage, error) {
	var images []*models.ProductImage
	err := r.db.Where("product_id = ?", productID).Find(&images).Error
	return images, err
}
