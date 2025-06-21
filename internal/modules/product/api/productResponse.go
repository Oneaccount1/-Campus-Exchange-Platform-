package api

import (
	"campus/internal/models"
	"time"
)

type ProductResponse struct {
	ID          uint                  `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	Images      []models.ProductImage `json:"images"`
	Category    string                `json:"category"`
	Condition   string                `json:"condition"`
	UserID      uint                  `json:"user_id"`
	Status      string                `json:"status"`
	SoldAt      time.Time             `json:"sold_at"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ProductListResponse struct {
	Products []*ProductResponse `json:"products"`
	Total    uint               `json:"total"`
	Page     uint               `json:"page"`
	Size     uint               `json:"size"`
}

// LatestProductItem 最新商品项
type LatestProductItem struct {
	ID         uint      `json:"id"`         // 商品ID
	Title      string    `json:"title"`      // 商品标题
	Price      float64   `json:"price"`      // 商品价格
	Category   string    `json:"category"`   // 商品分类
	Seller     string    `json:"seller"`     // 卖家名称
	CreateTime time.Time `json:"createTime"` // 创建时间
	Image      string    `json:"image"`      // 商品图片
	Status     string    `json:"status"`     // 商品状态
}

func ConvertToProductResponse(product *models.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Images:      product.ProductImages,
		Category:    product.Category,
		Condition:   product.Condition,
		UserID:      product.UserID,
		Status:      product.Status,
		SoldAt:      product.SoldAt,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ConvertToProductListResponse(products []*models.Product, total, page, size uint) *ProductListResponse {
	var productResponses []*ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ConvertToProductResponse(product))
	}
	return &ProductListResponse{
		Products: productResponses,
		Total:    total,
		Page:     page,
		Size:     size,
	}
}
