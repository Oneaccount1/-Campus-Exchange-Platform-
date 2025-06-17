package services

import (
	"campus/internal/models"
	"campus/internal/modules/order/api"
	"campus/internal/modules/order/repositories"
	"campus/internal/utils/errors"
)

type OrderService interface {
	CreateOrder(data *api.CreateOrderRequest) (*api.OrderResponse, error)
	DeleteOrder(id uint) error
	UpdateOrderStatus(id uint, data *api.UpdateOrderStatusRequest) (*api.OrderResponse, error)
	GetOrderByID(id uint) (*api.OrderResponse, error)
	GetUserOrders(buyerID uint, page, size uint) (*api.OrderListResponse, error)
}

type OrderServiceImpl struct {
	repository repositories.OrderRepository
}

func NewOrderService(orderRep repositories.OrderRepository) OrderService {
	return &OrderServiceImpl{
		repository: orderRep,
	}
}

func (s *OrderServiceImpl) CreateOrder(data *api.CreateOrderRequest) (*api.OrderResponse, error) {
	order := &models.Order{
		BuyerID:   data.BuyerID,
		SellerID:  data.SellerID,
		ProductID: data.ProductID,
		Status:    "卖家未处理",
	}

	if err := s.repository.Create(order); err != nil {
		return nil, errors.NewInternalServerError("创建订单失败", err)
	}

	return api.ConvertToOrderResponse(order), nil
}

func (s *OrderServiceImpl) DeleteOrder(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return errors.NewInternalServerError("删除订单失败", err)
	}
	return nil
}

func (s *OrderServiceImpl) UpdateOrderStatus(id uint, data *api.UpdateOrderStatusRequest) (*api.OrderResponse, error) {
	// 先获取订单当前状态
	order, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("订单", err)
	}

	// 校验订单当前状态是否为未处理，只有未处理状态才能修改为已同意或已拒绝
	if order.Status != "卖家未处理" {
		return nil, errors.NewBadRequestError("订单状态不可修改", nil)
	}

	if err := s.repository.UpdateStatus(id, data.Status); err != nil {
		return nil, errors.NewInternalServerError("更新订单状态失败", err)
	}

	// 重新获取更新后的订单信息
	updatedOrder, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("订单", err)
	}

	return api.ConvertToOrderResponse(updatedOrder), nil
}

func (s *OrderServiceImpl) GetOrderByID(id uint) (*api.OrderResponse, error) {
	order, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("订单", err)
	}

	return api.ConvertToOrderResponse(order), nil
}

func (s *OrderServiceImpl) GetUserOrders(buyerID uint, page, size uint) (*api.OrderListResponse, error) {
	orders, total, err := s.repository.GetByBuyerID(buyerID, page, size)
	if err != nil {
		return nil, errors.NewInternalServerError("查询用户订单失败", err)
	}

	return api.ConvertToOrderListResponse(orders, uint(total), page, size), nil
}
