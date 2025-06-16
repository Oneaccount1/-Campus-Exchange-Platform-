package api

import (
	"campus/internal/models"
	"campus/internal/modules/product/api"
	"time"
)

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	ID        uint                `json:"id"`
	ProductID uint                `json:"product_id"`
	Product   api.ProductResponse `json:"product"`
	CreatedAt time.Time           `json:"created_at"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	Favorites []FavoriteResponse `json:"favorites"`
	Total     int64              `json:"total"`
	Page      uint               `json:"page"`
	Size      uint               `json:"size"`
}

// ConvertToFavoriteResponse 将数据库模型转换为响应对象
func ConvertToFavoriteResponse(favorite *models.Favorite) FavoriteResponse {
	return FavoriteResponse{
		ID:        favorite.ID,
		ProductID: favorite.ProductID,
		Product:   *api.ConvertToProductResponse(&favorite.Product),
		CreatedAt: favorite.CreatedAt,
	}
}

// ConvertToFavoriteListResponse 将收藏列表转换为响应对象
func ConvertToFavoriteListResponse(favorites []*models.Favorite, total int64, page, size uint) FavoriteListResponse {
	var responses []FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, ConvertToFavoriteResponse(favorite))
	}

	return FavoriteListResponse{
		Favorites: responses,
		Total:     total,
		Page:      page,
		Size:      size,
	}
}
