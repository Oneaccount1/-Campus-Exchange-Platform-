package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"errors"
	"gorm.io/gorm"
)

// ProductRepository 商品数据访问接口
type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id string) (*models.Product, error)
	Create(product models.Product) (*models.Product, error)
	Update(id string, product models.Product) (*models.Product, error)
	Delete(id string) error
	SearchProductsByKeyword(keyword string) ([]models.Product, error)
}

// ProductRepositoryImpl 商品数据访问实现结构体
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository 创建新的商品数据访问实例
func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{
		db: bootstrap.GetDB(),
	}
}

// GetAll 获取所有商品
func (r *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

// GetByID 根据ID获取商品
func (r *ProductRepositoryImpl) GetByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Create 创建新商品
func (r *ProductRepositoryImpl) Create(product models.Product) (*models.Product, error) {
	err := r.db.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新商品信息
func (r *ProductRepositoryImpl) Update(id string, product models.Product) (*models.Product, error) {
	var existingProduct models.Product
	err := r.db.First(&existingProduct, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&existingProduct).Updates(product).Error
	if err != nil {
		return nil, err
	}

	return &existingProduct, nil
}

// Delete 删除商品
func (r *ProductRepositoryImpl) Delete(id string) error {
	err := r.db.Delete(&models.Product{}, "id = ?", id).Error
	return err
}

// SearchProductsByKeyword 通过商品名称模糊查询商品
func (r *ProductRepositoryImpl) SearchProductsByKeyword(keyword string) ([]models.Product, error) {
	var products []models.Product
	query := "%" + keyword + "%"
	err := r.db.Where("title LIKE ?", query).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
