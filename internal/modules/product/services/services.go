package services

import (
	"campus/internal/models"
	"campus/internal/modules/product/repositories"
)

// ProductService 商品服务接口
type ProductService interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id string) (*models.Product, error)
	CreateProduct(product models.Product) (*models.Product, error)
	UpdateProduct(id string, product models.Product) (*models.Product, error)
	DeleteProduct(id string) error
}

// ProductServiceImpl 商品服务实现结构体
type ProductServiceImpl struct {
	repository repositories.ProductRepository
}

// NewProductService 创建新的商品服务实例
func NewProductService() ProductService {
	return &ProductServiceImpl{
		repository: repositories.NewProductRepository(),
	}
}

// GetAllProducts 获取所有商品
func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.repository.GetAll()
}

// GetProductByID 根据ID获取商品
func (s *ProductServiceImpl) GetProductByID(id string) (*models.Product, error) {
	return s.repository.GetByID(id)
}

// CreateProduct 创建新商品
func (s *ProductServiceImpl) CreateProduct(product models.Product) (*models.Product, error) {
	return s.repository.Create(product)
}

// UpdateProduct 更新商品信息
func (s *ProductServiceImpl) UpdateProduct(id string, product models.Product) (*models.Product, error) {
	return s.repository.Update(id, product)
}

// DeleteProduct 删除商品
func (s *ProductServiceImpl) DeleteProduct(id string) error {
	return s.repository.Delete(id)
}
