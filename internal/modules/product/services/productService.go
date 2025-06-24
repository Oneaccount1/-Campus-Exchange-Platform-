package services

import (
	"campus/internal/models"
	"campus/internal/modules/product/api"
	"campus/internal/modules/product/repositories"
	"campus/internal/utils/errors"
	"fmt"
	"strconv"
	"time"
)

const (
	ProductStatusAvailable = "售卖中"
	ProductStatusSold      = "已下架"
	ProductStatusReviewing = "审核中"
)

type ProductService interface {
	GetAllProducts(page, size uint) (*api.ProductListResponse, error)
	GetProductByID(id string) (*api.ProductResponse, error)
	CreateProduct(data *api.CreateProductRequest) (*api.ProductResponse, error)
	UpdateProduct(id string, data *api.UpdateProductRequest) (*api.ProductResponse, error)
	DeleteProduct(id string) error
	SearchProductsByKeyword(keyword string, page, size uint) (*api.ProductListResponse, error)
	GetUserProducts(userID uint, page, size uint) (*api.ProductListResponse, error)
	GetSolvingProducts(page, size uint) (*api.ProductListResponse, error)
	FilterProducts(filter *api.FilterProductsRequest) (*api.ProductListResponse, error)
	GetLatestProducts(limit uint) (*api.ProductListResponse, error)
	UpdateProductStatus(id uint, status string) (*models.Product, error)
}

type ProductServiceImpl struct {
	productRep repositories.ProductRepository
	imageRep   repositories.ProductImageRepository
}

func (s *ProductServiceImpl) FilterProducts(filter *api.FilterProductsRequest) (*api.ProductListResponse, error) {
	products, tot, err := s.productRep.FilterProducts(filter)
	if err != nil {
		return nil, errors.NewInternalServerError("赛选商品失败", err)
	}
	return api.ConvertToProductListResponse(products, uint(tot), filter.Page, filter.Size), nil

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

func (s *ProductServiceImpl) GetSolvingProducts(page, size uint) (*api.ProductListResponse, error) {
	products, total, err := s.productRep.GetSolvingProducts(page, size)
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
		Status:      ProductStatusReviewing,
		SoldAt:      time.Now(),
	}

	productID, err := s.productRep.Create(product)
	if err != nil {
		return nil, errors.NewInternalServerError("创建商品失败", err)
	}

	// 添加对应的图片
	var images []models.ProductImage
	for _, imageURL := range data.Images {
		image := &models.ProductImage{
			ProductID: productID,
			ImageURL:  imageURL,
		}
		images = append(images, *image)
	}

	err = s.imageRep.CreateImages(images)
	if err != nil {
		return nil, err
	}
	product.ProductImages = images

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
		SoldAt:      time.Now(),
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
	var images []models.ProductImage
	pid, err := strconv.Atoi(id)
	for _, imageURL := range data.Images {
		image := &models.ProductImage{
			ProductID: uint(pid),
			ImageURL:  imageURL,
		}
		images = append(images, *image)
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

func (s *ProductServiceImpl) GetUserProducts(userID uint, page, size uint) (*api.ProductListResponse, error) {
	products, total, err := s.productRep.GetByUserID(userID, page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("查询用户发布商品失败", err)
	}

	return api.ConvertToProductListResponse(products, uint(total), page, size), nil
}

// GetLatestProducts 获取最新商品
func (s *ProductServiceImpl) GetLatestProducts(limit uint) (*api.ProductListResponse, error) {
	if limit == 0 {
		limit = 8 // 默认获取8条
	}

	products, total, err := s.productRep.GetLatest(limit)
	if err != nil {
		return nil, errors.NewInternalServerError("获取最新商品失败", err)
	}

	return api.ConvertToProductListResponse(products, uint(total), 1, limit), nil
}

// UpdateProductStatus 更新商品状态
func (s *ProductServiceImpl) UpdateProductStatus(productID uint, status string) (*models.Product, error) {
	product, err := s.productRep.GetByID(strconv.Itoa(int(productID)))
	if err != nil {
		return nil, errors.NewNotFoundError("找不到此商品", err)
	}
	fmt.Println(status)
	product.Status = status
	//fmt.Println(product)
	if err := s.productRep.Update(strconv.Itoa(int(productID)), product); err != nil {
		return nil, errors.NewInternalServerError("更新商品状态失败", err)
	}
	return product, nil
}
