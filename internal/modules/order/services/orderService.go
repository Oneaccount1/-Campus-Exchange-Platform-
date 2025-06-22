package services

import (
	"campus/internal/models"
	"campus/internal/modules/order/api"
	"campus/internal/modules/order/repositories"
	"campus/internal/utils/errors"
	"fmt"
	"time"
)

type OrderService interface {
	CreateOrder(data *api.CreateOrderRequest) (*api.OrderResponse, error)
	DeleteOrder(id uint) error
	UpdateOrderStatus(id uint, data *api.UpdateOrderStatusRequest) (*api.OrderResponse, error)
	GetOrderByID(id uint) (*api.OrderResponse, error)
	GetUserOrders(buyerID uint, page, size uint) (*api.OrderListResponse, error)

	// 管理员接口
	GetAdminOrderList(req *api.AdminOrderListRequest) (*api.AdminOrderListResponse, error)
	GetAdminOrderDetail(id uint) (*api.AdminOrderDetailResponse, error)
	AdminUpdateOrderStatus(id uint, req *api.AdminUpdateOrderStatusRequest) error
	ExportOrders(req *api.AdminOrderListRequest) ([]byte, error)
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

	if err := s.repository.UpdateStatus(id, data.Status, data.Remark); err != nil {
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

// GetAdminOrderList 管理员获取订单列表
func (s *OrderServiceImpl) GetAdminOrderList(req *api.AdminOrderListRequest) (*api.AdminOrderListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 获取订单列表
	orders, total, err := s.repository.GetOrdersForAdmin(
		req.Search, req.Status, req.StartDate, req.EndDate, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.NewInternalServerError("获取订单列表失败", err)
	}

	// 构建响应
	response := &api.AdminOrderListResponse{
		Total: total,
		List:  make([]api.AdminOrderItem, 0, len(orders)),
	}

	// 填充订单数据
	for _, order := range orders {
		// 获取商品信息
		var productTitle string
		var productImage string

		if order.Product.ID > 0 {
			productTitle = order.Product.Title
			if len(order.Product.ProductImages) > 0 {
				productImage = order.Product.ProductImages[0].ImageURL
			}
		}

		// 获取买家和卖家信息
		var buyerName string
		var sellerName string

		if order.Buyer.ID > 0 {
			buyerName = order.Buyer.Username
		}

		if order.Seller.ID > 0 {
			sellerName = order.Seller.Username
		}

		// 构建订单项
		item := api.AdminOrderItem{
			ID:           fmt.Sprintf("O%d", order.ID),
			ProductTitle: productTitle,
			ProductImage: productImage,
			Buyer:        buyerName,
			Seller:       sellerName,
			Status:       order.Status,
			CreateTime:   order.CreatedAt,
		}

		// 添加支付时间和完成时间（如果有）
		if !order.PayTime.IsZero() {
			item.PayTime = order.PayTime
		}

		if !order.CompleteTime.IsZero() {
			item.CompleteTime = order.CompleteTime
		}

		response.List = append(response.List, item)
	}

	return response, nil
}

// GetAdminOrderDetail 管理员获取订单详情
func (s *OrderServiceImpl) GetAdminOrderDetail(id uint) (*api.AdminOrderDetailResponse, error) {
	// 获取订单详情
	order, err := s.repository.GetOrderDetailForAdmin(id)
	if err != nil {
		return nil, errors.NewNotFoundError("订单", err)
	}

	// 获取订单日志
	logs, err := s.repository.GetOrderLogs(id)
	if err != nil {
		return nil, errors.NewInternalServerError("获取订单日志失败", err)
	}

	// 构建响应
	response := &api.AdminOrderDetailResponse{
		ID:        fmt.Sprintf("O%d", order.ID),
		ProductID: order.ProductID,
		
		BuyerID:    order.BuyerID,
		SellerID:   order.SellerID,
		Status:     order.Status,
		CreateTime: order.CreatedAt,
		Remark:     order.Remark,
		Logs:       make([]api.OrderLogItem, 0, len(logs)),
	}

	// 填充商品信息
	if order.Product.ID > 0 {
		response.ProductTitle = order.Product.Title
		if len(order.Product.ProductImages) > 0 {
			response.ProductImage = order.Product.ProductImages[0].ImageURL
		}
	}

	// 填充买家信息
	if order.Buyer.ID > 0 {
		response.Buyer = order.Buyer.Username
		response.BuyerPhone = order.Buyer.Phone
		// 这里假设买家地址存储在描述字段中，实际应用中可能需要单独的地址字段
		response.BuyerAddress = order.Buyer.Description
	}

	// 填充卖家信息
	if order.Seller.ID > 0 {
		response.Seller = order.Seller.Username
		response.SellerPhone = order.Seller.Phone
	}

	// 填充时间信息
	if !order.PayTime.IsZero() {
		response.PayTime = order.PayTime
	}

	if !order.DeliveryTime.IsZero() {
		response.DeliveryTime = order.DeliveryTime
	}

	if !order.CompleteTime.IsZero() {
		response.CompleteTime = order.CompleteTime
	}

	// 填充日志信息
	for _, log := range logs {
		logItem := api.OrderLogItem{
			Action:   log.Action,
			Time:     log.CreatedAt,
			Operator: log.Operator,
			Remark:   log.Remark,
		}
		response.Logs = append(response.Logs, logItem)
	}

	return response, nil
}

// AdminUpdateOrderStatus 管理员更新订单状态
func (s *OrderServiceImpl) AdminUpdateOrderStatus(id uint, req *api.AdminUpdateOrderStatusRequest) error {
	// 获取订单
	order, err := s.repository.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("订单", err)
	}

	// 根据状态更新相关时间字段
	now := time.Now()
	switch req.Status {
	case "待付款":
		// 无需特殊处理
	case "待发货":
		order.PayTime = now
	case "待收货":
		order.DeliveryTime = now
	case "已完成":
		order.CompleteTime = now
	case "已取消":
		// 无需特殊处理
	}

	// 更新订单状态和时间
	if err := s.repository.UpdateStatus(id, req.Status, req.Remark); err != nil {
		return errors.NewInternalServerError("更新订单状态失败", err)
	}

	return nil
}

// ExportOrders 导出订单数据
func (s *OrderServiceImpl) ExportOrders(req *api.AdminOrderListRequest) ([]byte, error) {
	// 设置导出时的页码和每页数量
	req.Page = 1
	req.PageSize = 1000 // 导出更多数据

	// 获取订单列表
	result, err := s.GetAdminOrderList(req)
	if err != nil {
		return nil, err
	}

	// 构建CSV数据
	csvData := "订单ID,商品标题,价格,买家,卖家,状态,创建时间,支付时间,完成时间\n"

	for _, order := range result.List {
		// 格式化时间
		createTime := order.CreateTime.Format("2006-01-02 15:04:05")
		payTime := ""
		if !order.PayTime.IsZero() {
			payTime = order.PayTime.Format("2006-01-02 15:04:05")
		}
		completeTime := ""
		if !order.CompleteTime.IsZero() {
			completeTime = order.CompleteTime.Format("2006-01-02 15:04:05")
		}

		// 添加一行数据
		line := fmt.Sprintf("%s,%s,%.2f,%s,%s,%s,%s,%s,%s\n",
			order.ID,
			order.ProductTitle,
			order.Price,
			order.Buyer,
			order.Seller,
			order.Status,
			createTime,
			payTime,
			completeTime,
		)

		csvData += line
	}

	return []byte(csvData), nil
}
