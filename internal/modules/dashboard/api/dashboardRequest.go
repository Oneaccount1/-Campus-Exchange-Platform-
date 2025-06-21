package api

// ProductTrendRequest 商品发布趋势请求
type ProductTrendRequest struct {
	TimeRange string `form:"timeRange" json:"timeRange"` // 时间范围，可选值：week, month
}

// LatestProductsRequest 最新商品列表请求
type LatestProductsRequest struct {
	Limit int `form:"limit" json:"limit"` // 限制数量
}

// ActivitiesRequest 系统活动请求
type ActivitiesRequest struct {
	Limit int `form:"limit" json:"limit"` // 限制数量
} 