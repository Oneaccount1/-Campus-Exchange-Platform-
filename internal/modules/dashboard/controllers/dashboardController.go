package controllers

import (
	"campus/internal/modules/dashboard/api"
	"campus/internal/modules/dashboard/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	service services.DashboardService
}

// NewDashboardController 创建仪表盘控制器实例
func NewDashboardController() *DashboardController {
	return &DashboardController{
		service: services.NewDashboardService(),
	}
}

// GetStats 获取仪表盘统计数据
func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.service.GetStats()
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "获取成功", stats)
}

// GetProductTrend 获取商品发布趋势
func (c *DashboardController) GetProductTrend(ctx *gin.Context) {
	var req api.ProductTrendRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 使用timeRange参数，默认为week
	timeRange := req.TimeRange
	if timeRange != "week" && timeRange != "month" {
		timeRange = "week"
	}

	trend, err := c.service.GetProductTrend(timeRange)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "获取成功", trend)
}

// GetCategoryStats 获取商品分类统计
func (c *DashboardController) GetCategoryStats(ctx *gin.Context) {
	stats, err := c.service.GetCategoryStats()
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "获取成功", stats)
}

// GetLatestProducts 获取最新商品
func (c *DashboardController) GetLatestProducts(ctx *gin.Context) {
	var req api.LatestProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 使用limit参数，默认为5
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}

	products, err := c.service.GetLatestProducts(limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "获取成功", products)
}

// GetActivities 获取系统活动
func (c *DashboardController) GetActivities(ctx *gin.Context) {
	var req api.ActivitiesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	// 使用limit参数，默认为5
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}

	activities, err := c.service.GetActivities(limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "获取成功", activities)
}
