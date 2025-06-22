package controllers

import (
	"campus/internal/modules/order/api"
	"campus/internal/modules/order/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderController struct {
	service services.OrderService
}

func NewOrderController(srv services.OrderService) *OrderController {
	return &OrderController{
		service: srv,
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req api.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	order, err := c.service.CreateOrder(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "订单创建成功", order)
}

func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效订单ID", err))
		return
	}

	if err := c.service.DeleteOrder(uint(id)); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "订单删除成功", nil)
}

func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效订单ID", err))
		return
	}

	var req api.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	order, err := c.service.UpdateOrderStatus(uint(id), &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "订单状态更新成功", order)
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效订单ID", err))
		return
	}

	order, err := c.service.GetOrderByID(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, order)
}

func (c *OrderController) GetUserOrders(ctx *gin.Context) {
	var req api.GetUserOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	orders, err := c.service.GetUserOrders(req.UserID, req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, orders)
}

// GetAdminOrderList 管理员获取订单列表
func (c *OrderController) GetAdminOrderList(ctx *gin.Context) {
	var req api.AdminOrderListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	
	// 获取订单列表
	result, err := c.service.GetAdminOrderList(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	
	response.SuccessWithMessage(ctx, "获取成功", result)
}

// GetAdminOrderDetail 管理员获取订单详情
func (c *OrderController) GetAdminOrderDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效订单ID", err))
		return
	}
	
	// 获取订单详情
	result, err := c.service.GetAdminOrderDetail(uint(id))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	
	response.SuccessWithMessage(ctx, "获取成功", result)
}

// AdminUpdateOrderStatus 管理员更新订单状态
func (c *OrderController) AdminUpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效订单ID", err))
		return
	}
	
	var req api.AdminUpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	
	// 更新订单状态
	if err := c.service.AdminUpdateOrderStatus(uint(id), &req); err != nil {
		response.HandleError(ctx, err)
		return
	}
	
	response.SuccessWithMessage(ctx, "更新成功", nil)
}

// ExportOrders 导出订单数据
func (c *OrderController) ExportOrders(ctx *gin.Context) {
	var req api.AdminOrderListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	
	// 导出订单数据
	data, err := c.service.ExportOrders(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	
	// 设置响应头
	ctx.Header("Content-Disposition", "attachment; filename=orders.csv")
	ctx.Header("Content-Type", "text/csv")
	ctx.Data(http.StatusOK, "text/csv", data)
}
