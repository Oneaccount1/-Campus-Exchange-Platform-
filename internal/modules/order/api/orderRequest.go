package api

type CreateOrderRequest struct {
	BuyerID   uint    `json:"buyer_id" binding:"required"`
	SellerID  uint    `json:"seller_id" binding:"required"`
	ProductID uint    `json:"product_id" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	// 明确状态可选值为卖家已同意或卖家已拒绝
	Status string `json:"status" binding:"required,oneof=卖家已同意 卖家已拒绝"`
}

type GetUserOrdersRequest struct {
	// 这里可根据需求添加分页等参数，当前无需额外参数
	Page uint `json:"page" binding:"required"`
	Size uint `json:"size" binding:"required"`
}
