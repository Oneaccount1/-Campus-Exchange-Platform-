package api

import (
	"campus/internal/models"
	"time"
)

type OrderResponse struct {
	ID           uint       `json:"id"`
	BuyerID      uint       `json:"buyer_id"`
	SellerID     uint       `json:"seller_id"`
	ProductID    uint       `json:"product_id"`
	Status       string     `json:"status"`
	Price        float64    `json:"price"`
	PayTime      *time.Time `json:"pay_time,omitempty"`
	DeliveryTime *time.Time `json:"delivery_time,omitempty"`
	CompleteTime *time.Time `json:"complete_time,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func ConvertToOrderResponse(order *models.Order) *OrderResponse {
	return &OrderResponse{
		ID:           order.ID,
		BuyerID:      order.BuyerID,
		SellerID:     order.SellerID,
		ProductID:    order.ProductID,
		Status:       order.Status,
		PayTime:      order.PayTime,
		DeliveryTime: order.DeliveryTime,
		CompleteTime: order.CompleteTime,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
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

// AdminOrderItem 管理员订单列表项
type AdminOrderItem struct {
	ID           string     `json:"id"`
	ProductTitle string     `json:"productTitle"`
	ProductImage string     `json:"productImage"`
	Price        float64    `json:"price"`
	Buyer        string     `json:"buyer"`
	Seller       string     `json:"seller"`
	Status       string     `json:"status"`
	CreateTime   time.Time  `json:"createTime"`
	PayTime      *time.Time `json:"payTime,omitempty"`
	CompleteTime *time.Time `json:"completeTime,omitempty"`
}

// AdminOrderListResponse 管理员订单列表响应
type AdminOrderListResponse struct {
	Total int64            `json:"total"`
	List  []AdminOrderItem `json:"list"`
}

// OrderLogItem 订单日志项
type OrderLogItem struct {
	Action   string    `json:"action"`
	Time     time.Time `json:"time"`
	Operator string    `json:"operator"`
	Remark   string    `json:"remark,omitempty"`
}

// AdminOrderDetailResponse 管理员订单详情响应
type AdminOrderDetailResponse struct {
	ID           string         `json:"id"`
	ProductID    uint           `json:"productId"`
	ProductTitle string         `json:"productTitle"`
	ProductImage string         `json:"productImage"`
	Price        float64        `json:"price"`
	BuyerID      uint           `json:"buyerId"`
	Buyer        string         `json:"buyer"`
	BuyerPhone   string         `json:"buyerPhone"`
	BuyerAddress string         `json:"buyerAddress"`
	SellerID     uint           `json:"sellerId"`
	Seller       string         `json:"seller"`
	SellerPhone  string         `json:"sellerPhone"`
	Status       string         `json:"status"`
	CreateTime   time.Time      `json:"createTime"`
	PayTime      *time.Time     `json:"payTime,omitempty"`
	DeliveryTime *time.Time     `json:"deliveryTime,omitempty"`
	CompleteTime *time.Time     `json:"completeTime,omitempty"`
	Remark       string         `json:"remark,omitempty"`
	Logs         []OrderLogItem `json:"logs"`
}
