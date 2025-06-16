package api

// FavoriteRequest 收藏商品的请求
type FavoriteRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// QueryPageRequest 分页查询请求参数
type QueryPageRequest struct {
	Page uint `form:"page" binding:"required,min=1"`
	Size uint `form:"size" binding:"required,min=1,max=100"`
}
