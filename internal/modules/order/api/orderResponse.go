package api

import (
	"campus/internal/models"
	"time"
)

type OrderResponse struct {
	ID        uint      `json:"id"`
	BuyerID   uint      `json:"buyer_id"`
	SellerID  uint      `json:"seller_id"`
	ProductID uint      `json:"product_id"`
	Status    string    `json:"status"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ConvertToOrderResponse(order *models.Order) *OrderResponse {
	return &OrderResponse{
		ID:        order.ID,
		BuyerID:   order.BuyerID,
		SellerID:  order.SellerID,
		ProductID: order.ProductID,
		Status:    order.Status,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  uint             `json:"total"`
	Page   uint             `json:"page"`
	Size   uint             `json:"size"`
}

func ConvertToOrderListResponse(orders []*models.Order, total, page, size uint) *OrderListResponse {
	var orderResponses []*OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, ConvertToOrderResponse(order))
	}
	return &OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Size:   size,
	}
}
