package services

import (
	"campus/internal/models"
	"campus/internal/modules/product/api"
	"campus/internal/modules/product/repositories"
	"campus/internal/utils/errors"
	"strconv"
)

type ProductService interface {
	GetAllProducts(page, size uint) (*api.ProductListResponse, error)
	GetProductByID(id string) (*api.ProductResponse, error)
	CreateProduct(data *api.CreateProductRequest) (*api.ProductResponse, error)
	UpdateProduct(id string, data *api.UpdateProductRequest) (*api.ProductResponse, error)
	DeleteProduct(id string) error
	SearchProductsByKeyword(keyword string, page, size uint) (*api.ProductListResponse, error)
}

type ProductServiceImpl struct {
	productRep repositories.ProductRepository
	imageRep   repositories.ProductImageRepository
}

func NewProductService() ProductService {
	return &ProductServiceImpl{
		productRep: repositories.NewProductRepository(),
		imageRep:   repositories.NewProductImageRepository(),
	}
}

func (s *ProductServiceImpl) GetAllProducts(page, size uint) (*api.ProductListResponse, error) {
	products, total, err := s.productRep.GetAll(page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("查询商品列表失败", err)
	}

	return api.ConvertToProductListResponse(products, uint(total), page, size), nil
}

func (s *ProductServiceImpl) GetProductByID(id string) (*api.ProductResponse, error) {
	product, err := s.productRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("商品", err)
	}

	return api.ConvertToProductResponse(product), nil
}

func (s *ProductServiceImpl) CreateProduct(data *api.CreateProductRequest) (*api.ProductResponse, error) {
	product := &models.Product{
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		Category:    data.Category,
		Condition:   data.Condition,
		UserID:      data.UserID,
		Status:      data.Status,
		SoldAt:      data.SoldAt,
	}

	productID, err := s.productRep.Create(product)
	if err != nil {
		return nil, errors.NewInternalServerError("创建商品失败", err)
	}

	// 添加对应的图片
	var images []*models.ProductImage
	for _, imageURL := range data.Images {
		image := &models.ProductImage{
			ProductID: productID,
			ImageURL:  imageURL,
		}
		images = append(images, image)
	}

	err = s.imageRep.CreateImages(images)
	if err != nil {
		return nil, err
	}

	return api.ConvertToProductResponse(product), nil
}

func (s *ProductServiceImpl) UpdateProduct(id string, data *api.UpdateProductRequest) (*api.ProductResponse, error) {
	_, err := s.productRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("商品", err)
	}

	updatedProduct := &models.Product{
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		Category:    data.Category,
		Condition:   data.Condition,
		Status:      data.Status,
		SoldAt:      data.SoldAt,
	}

	if err := s.productRep.Update(id, updatedProduct); err != nil {
		return nil, errors.NewInternalServerError("更新商品失败", err)
	}

	/*
		更新图片
	*/

	// 删除旧图片
	if err := s.imageRep.DeleteImagesByProductID(id); err != nil {
		return nil, errors.NewInternalServerError("删除旧图片失败", err)
	}

	// 添加新图片
	var images []*models.ProductImage
	pid, err := strconv.Atoi(id)
	for _, imageURL := range data.Images {
		image := &models.ProductImage{
			ProductID: uint(pid),
			ImageURL:  imageURL,
		}
		images = append(images, image)
	}

	if err := s.imageRep.CreateImages(images); err != nil {
		return nil, errors.NewInternalServerError("添加新图片失败", err)
	}

	updated, err := s.productRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("商品", err)
	}

	return api.ConvertToProductResponse(updated), nil
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	if err := s.productRep.Delete(id); err != nil {
		return errors.NewInternalServerError("删除商品失败", err)
	}
	return nil
}

func (s *ProductServiceImpl) SearchProductsByKeyword(keyword string, page, size uint) (*api.ProductListResponse, error) {
	products, total, err := s.productRep.SearchProductsByKeyword(keyword, page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("搜索商品失败", err)
	}

	return api.ConvertToProductListResponse(products, uint(total), page, size), nil
}
