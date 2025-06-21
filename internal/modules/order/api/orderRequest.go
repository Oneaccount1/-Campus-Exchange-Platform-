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
	Remark string `json:"remark"` // 备注信息，可选
}

type GetUserOrdersRequest struct {
	UserID uint `json:"user_id" form:"user_id" binding:"required"`
	Page   uint `json:"page" form:"page" binding:"required,min=1"`
	Size   uint `json:"size" form:"size" binding:"required,min=1,max=100"`
}

// AdminOrderListRequest 管理员获取订单列表请求
type AdminOrderListRequest struct {
	Page      uint   `json:"page" form:"page"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
	Search    string `json:"search" form:"search"`
	Status    string `json:"status" form:"status"`
	StartDate string `json:"startDate" form:"startDate"`
	EndDate   string `json:"endDate" form:"endDate"`
}

// AdminUpdateOrderStatusRequest 管理员更新订单状态请求
type AdminUpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Remark string `json:"remark"` // 备注信息，可选
}
