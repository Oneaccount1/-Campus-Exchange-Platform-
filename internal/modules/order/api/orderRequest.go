package api

type CreateOrderRequest struct {
	BuyerID   uint `json:"buyer_id" binding:"required"`
	SellerID  uint `json:"seller_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
	//Price     float64 `json:"price" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	// 明确状态可选值为卖家已同意或卖家已拒绝
	Status string `json:"status" binding:"required,oneof=卖家已同意 卖家已拒绝"`
}

type GetUserOrdersRequest struct {
	UserID uint `json:"user_id" form:"user_id" binding:"required"`
	Page   uint `json:"page" form:"page" binding:"required,min=1"`
	Size   uint `json:"size" form:"size" binding:"required,min=1,max=100"`
}
