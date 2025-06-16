package services

import (
	"campus/internal/models"
	api2 "campus/internal/modules/product/api"
	prodRep "campus/internal/modules/product/repositories"
	"campus/internal/modules/user/api"
	userRep "campus/internal/modules/user/repositories"
	"campus/internal/utils/errors"
)

type FavoriteService interface {
	// AddFavorite 添加收藏
	AddFavorite(userID uint, productID uint) error
	// RemoveFavorite 取消收藏
	RemoveFavorite(userID uint, productID uint) error
	// ListUserFavorites 获取用户收藏列表
	ListUserFavorites(userID uint, page, size uint) (*api.FavoriteListResponse, error)
	// CheckIsFavorite 检查是否已收藏
	CheckIsFavorite(userID uint, productID uint) (bool, error)
	// GetUserProducts 获取用户发布的商品
	GetUserProducts(userID uint, page, size uint) (*api2.ProductListResponse, error)
}

type favoriteService struct {
	favoriteRepo userRep.FavoriteRepository
	productRepo  prodRep.ProductRepository
}

func (f *favoriteService) AddFavorite(userID uint, productID uint) error {
	// 检查是否收藏
	isFavorite, err := f.favoriteRepo.CheckIsFavorite(userID, productID)
	if err != nil {
		return errors.NewInternalServerError("检查收藏状态失败", err)
	}
	if isFavorite {
		return errors.NewBadRequestError("已经收藏过该商品", nil)
	}
	// 创建收藏记录
	favorite := &models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	if err = f.favoriteRepo.Create(favorite); err != nil {
		return errors.NewInternalServerError("添加收藏失败", err)
	}
	return nil
}

// RemoveFavorite 取消收藏
func (f *favoriteService) RemoveFavorite(userID uint, productID uint) error {
	// 检查是否已经收藏过
	isFavorite, err := f.favoriteRepo.CheckIsFavorite(userID, productID)
	if err != nil {
		return errors.NewInternalServerError("检查收藏状态失败", err)
	}

	if !isFavorite {
		return errors.NewBadRequestError("未收藏该商品", nil)
	}
	// 删除收藏记录
	if err := f.favoriteRepo.Delete(userID, productID); err != nil {
		return errors.NewInternalServerError("取消收藏失败", err)
	}
	return nil

}

// ListUserFavorites 获取用户收藏列表
func (f *favoriteService) ListUserFavorites(userID uint, page, size uint) (*api.FavoriteListResponse, error) {
	favorites, tot, err := f.favoriteRepo.GetByUser(userID, page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("获取收藏列表失败", err)
	}
	response := api.ConvertToFavoriteListResponse(favorites, tot, page, size)
	return &response, err
}

func (f *favoriteService) CheckIsFavorite(userID uint, productID uint) (bool, error) {
	isFavorite, err := f.favoriteRepo.CheckIsFavorite(userID, productID)
	if err != nil {
		return false, errors.NewInternalServerError("检查收藏状态失败", err)
	}
	return isFavorite, nil
}

func (f *favoriteService) GetUserProducts(userID uint, page, size uint) (*api2.ProductListResponse, error) {
	products, tot, err := f.productRepo.GetProductsByUserID(userID, page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("获取用户发布的商品失败", err)
	}

	var response api2.ProductListResponse
	var productResponse []*api2.ProductResponse

	for _, product := range products {
		productResponse = append(productResponse, api2.ConvertToProductResponse(product))
	}
	response = api2.ProductListResponse{
		Products: productResponse,
		Total:    uint(tot),
		Page:     page,
		Size:     size,
	}
	return &response, nil
}

func NewFavoriteService() FavoriteService {
	return &favoriteService{
		favoriteRepo: userRep.NewFavoriteRepository(),
		productRepo:  prodRep.NewProductRepository(),
	}
}
