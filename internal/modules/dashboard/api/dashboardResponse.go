package api

import (
	"time"
)

// DashboardStatsResponse 仪表盘统计数据响应
type DashboardStatsResponse struct {
	ProductCount        int     `json:"productCount"`        // 商品总数
	ProductTrend        float64 `json:"productTrend"`        // 商品增长趋势，百分比
	UserCount           int     `json:"userCount"`           // 用户总数
	UserTrend           float64 `json:"userTrend"`           // 用户增长趋势，百分比
	NewUserCount        int     `json:"newUserCount"`        // 本周新增用户数
	TodayProductCount   int     `json:"todayProductCount"`   // 今日新增商品数
	TodayProductTrend   float64 `json:"todayProductTrend"`   // 今日商品增长趋势，百分比
	YesterdayProductCount int   `json:"yesterdayProductCount"` // 昨日新增商品数
	TotalAmount         float64 `json:"totalAmount"`         // 交易总额
	AmountTrend         float64 `json:"amountTrend"`         // 交易额增长趋势，百分比
	MonthAmount         float64 `json:"monthAmount"`         // 本月交易额
}

// ProductTrendItem 商品发布趋势项
type ProductTrendItem struct {
	Label string `json:"label"` // 标签，如"周一"
	Value int    `json:"value"` // 值，如商品数量
}

// ProductTrendResponse 商品发布趋势响应
type ProductTrendResponse []ProductTrendItem

// CategoryStatsItem 分类统计项
type CategoryStatsItem struct {
	Name  string `json:"name"`  // 分类名称
	Value int    `json:"value"` // 值，如商品数量
}

// CategoryStatsResponse 分类统计响应
type CategoryStatsResponse []CategoryStatsItem

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

// LatestProductsResponse 最新商品列表响应
type LatestProductsResponse []LatestProductItem

// ActivityItem 系统活动项
type ActivityItem struct {
	Content string `json:"content"` // 活动内容
	Time    string `json:"time"`    // 活动时间，如"10分钟前"
	Type    string `json:"type"`    // 活动类型
	Color   string `json:"color"`   // 活动颜色
}

// ActivitiesResponse 系统活动响应
type ActivitiesResponse []ActivityItem 