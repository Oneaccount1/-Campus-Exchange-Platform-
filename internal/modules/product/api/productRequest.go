package api

type CreateProductRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required"`
	Images      []string `json:"images"`
	Category    string   `json:"category"`
	Condition   string   `json:"condition"`
	UserID      uint     `json:"user_id" binding:"required"`
	Status      string   `json:"status" binding:"required,oneof=售卖中 已下架 审核中"`
}

type UpdateProductRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	Category    string   `json:"category"`
	Condition   string   `json:"condition"`
	Status      string   `json:"status" binding:"required,oneof=售卖中 已下架 审核中"`
}

type GetProductsRequest struct {
	Page uint `json:"page" form:"page" binding:"required,min=1"`
	Size uint `json:"size" form:"size" binding:"required,min=1,max=100"`
}

type GetUserProductsRequest struct {
	UserID uint `json:"user_id" form:"user_id" binding:"required"`
	Page   uint `json:"page" form:"page" binding:"required,min=1"`
	Size   uint `json:"size" form:"size" binding:"required,min=1,max=100"`
}

// BatchUpdateStatusRequest 批量更新商品状态请求
type BatchUpdateStatusRequest struct {
	ProductIDs []uint `json:"productIds" binding:"required,min=1"`         // 商品ID列表
	Status     string `json:"status" binding:"required,oneof=售卖中 已下架 待审核"` // 商品状态
}
type FilterProductsRequest struct {
	Page      uint   `json:"page" form:"page" binding:"required,min=1"`
	Size      uint   `json:"size" form:"size" binding:"required,min=1,max=100"`
	Keyword   string `json:"keyword" form:"keyword"`       // 关键词搜索
	Category  string `json:"category" form:"category"`     // 分类筛选
	Status    string `json:"status" form:"status"`         // 状态筛选
	Condition string `json:"condition" form:"condition"`   // 商品状况
	StartDate string `json:"start_date" form:"start_date"` // 开始日期
	EndDate   string `json:"end_date" form:"end_date"`     // 结束日期
}
