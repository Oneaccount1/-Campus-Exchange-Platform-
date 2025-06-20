package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(page, size uint) ([]*models.Product, int64, error)
	GetByID(id string) (*models.Product, error)
	Create(product *models.Product) (uint, error)
	Update(id string, product *models.Product) error
	Delete(id string) error
	SearchProductsByKeyword(keyword string, page, size uint) ([]*models.Product, int64, error)
	GetByUserID(userID uint, page, size uint) ([]*models.Product, int64, error)
	GetSolvingProducts(page, size uint) ([]*models.Product, int64, error)
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{
		db: bootstrap.GetDB(),
	}
}

func (r *ProductRepositoryImpl) GetAll(page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	err := r.db.Model(&models.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	// 使用 Preload 预加载图片
	err = r.db.Preload("ProductImages").Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetSolvingProducts(page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	err := r.db.Model(&models.Product{}).Where("status = ?", "售卖中").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	// 使用 Preload 预加载图片
	err = r.db.Preload("ProductImages").Where("status = ?", "售卖中").Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetByID(id string) (*models.Product, error) {
	var product models.Product
	// 使用 Preload 预加载图片
	err := r.db.Preload("ProductImages").First(&product, "id = ?", id).Error
	return &product, err
}

func (r *ProductRepositoryImpl) Create(product *models.Product) (uint, error) {
	err := r.db.Create(product).Error
	return product.ID, err
}

func (r *ProductRepositoryImpl) Update(id string, product *models.Product) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(product).Error
}

func (r *ProductRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryImpl) SearchProductsByKeyword(keyword string, page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	query := "%" + keyword + "%"

	err := r.db.Model(&models.Product{}).Where("title LIKE ?", query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.Where("title LIKE ?", query).Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetByUserID(userID uint, page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	// 计算总记录数
	err := r.db.Model(&models.Product{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err = r.db.Preload("ProductImages").Where("user_id = ?", userID).Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}
