package controllers

import (
	"campus/internal/modules/order/api"
	"campus/internal/modules/order/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
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
	// 假设从 JWT 中间件中获取当前用户 ID
	buyerID, exists := ctx.Get("userID")
	if !exists {
		response.HandleError(ctx, errors.NewUnauthorizedError("未获取到用户 ID", nil))
		return
	}

	// 类型转换
	userID, ok := buyerID.(uint)
	if !ok {
		response.HandleError(ctx, errors.NewInternalServerError("用户 ID 类型转换失败", nil))
		return
	}

	var req api.GetUserOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	orders, err := c.service.GetUserOrders(userID, req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, orders)
}
