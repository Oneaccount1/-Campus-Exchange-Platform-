package repositories

import (
	"campus/internal/models"
	"gorm.io/gorm"
	"time"
)

type OrderRepository interface {
	Create(order *models.Order) error
	Delete(id uint) error
	UpdateStatus(id uint, status string, remark string) error
	UpdateStatusWithTimes(id uint, status, remark string, payTime, deliveryTime, completeTime *time.Time) error
	GetByID(id uint) (*models.Order, error)
	GetByBuyerID(buyerID uint, page, size uint) ([]*models.Order, int64, error)
	// 管理员接口
	GetOrdersForAdmin(search, status, startDate, endDate string, page, pageSize uint) ([]*models.Order, int64, error)
	GetOrderDetailForAdmin(id uint) (*models.Order, error)
	CreateOrderLog(orderID uint, action, operator, remark string) error
	GetOrderLogs(orderID uint) ([]models.OrderLog, error)
	// 获取商品信息
	GetProductByID(productID uint) (*models.Product, error)
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

// GetProductByID 获取商品信息
func (r *OrderRepositoryImpl) GetProductByID(productID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, productID).Error
	return &product, err
}

func (r *OrderRepositoryImpl) Create(order *models.Order) error {
	err := r.db.Create(order).Error
	if err != nil {
		return err
	}

	// 创建订单日志
	return r.CreateOrderLog(order.ID, "创建订单", "用户", "")
}

func (r *OrderRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *OrderRepositoryImpl) UpdateStatus(id uint, status string, remark string) error {
	err := r.db.Model(&models.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"remark": remark,
	}).Error

	if err != nil {
		return err
	}

	// 创建订单状态变更日志
	var operator string
	if status == "卖家已同意" {
		operator = "卖家"
	} else if status == "卖家已拒绝" {
		operator = "卖家"
	} else {
		operator = "管理员"
	}

	return r.CreateOrderLog(id, "更新订单状态为"+status, operator, remark)
}

func (r *OrderRepositoryImpl) UpdateStatusWithTimes(id uint, status, remark string, payTime, deliveryTime, completeTime *time.Time) error {
	// 构建更新字段
	updates := map[string]interface{}{
		"status": status,
		"remark": remark,
	}

	// 只有当提供了时间才更新相应字段
	if payTime != nil {
		updates["pay_time"] = payTime
	}

	if deliveryTime != nil {
		updates["delivery_time"] = deliveryTime
	}

	if completeTime != nil {
		updates["complete_time"] = completeTime
	}

	// 更新订单
	err := r.db.Model(&models.Order{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}

	// 创建订单状态变更日志
	var operator string
	if status == "卖家已同意" {
		operator = "卖家"
	} else if status == "卖家已拒绝" {
		operator = "卖家"
	} else {
		operator = "管理员"
	}

	return r.CreateOrderLog(id, "更新订单状态为"+status, operator, remark)
}

func (r *OrderRepositoryImpl) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *OrderRepositoryImpl) GetByBuyerID(buyerID uint, page, size uint) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	// 计算总记录数
	err := r.db.Model(&models.Order{}).Where("buyer_id = ?", buyerID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err = r.db.Where("buyer_id = ?", buyerID).Offset(int(offset)).Limit(int(size)).Find(&orders).Error
	return orders, total, err
}

// GetOrdersForAdmin 管理员获取订单列表
func (r *OrderRepositoryImpl) GetOrdersForAdmin(search, status, startDate, endDate string, page, pageSize uint) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	// 添加搜索条件
	if search != "" {
		// 通过关联查询，搜索订单号、商品标题、买家名称、卖家名称
		query = query.Joins("JOIN products ON orders.product_id = products.id").
			Joins("JOIN users AS buyers ON orders.buyer_id = buyers.id").
			Joins("JOIN users AS sellers ON orders.seller_id = sellers.id").
			Where("orders.id LIKE ? OR products.title LIKE ? OR buyers.username LIKE ? OR sellers.username LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 添加状态筛选
	if status != "" {
		query = query.Where("orders.status = ?", status)
	}

	// 添加日期筛选
	if startDate != "" {
		startTime, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("orders.created_at >= ?", startTime)
		}
	}

	if endDate != "" {
		endTime, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// 增加一天，使得结束日期是包含当天的
			endTime = endTime.AddDate(0, 0, 1)
			query = query.Where("orders.created_at < ?", endTime)
		}
	}

	// 计算总记录数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询 - 添加预加载
	offset := (page - 1) * pageSize
	err = query.Preload("Product").
		Preload("Product.ProductImages").
		Preload("Buyer").
		Preload("Seller").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("orders.created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// GetOrderDetailForAdmin 管理员获取订单详情
func (r *OrderRepositoryImpl) GetOrderDetailForAdmin(id uint) (*models.Order, error) {
	var order models.Order

	// 使用预加载获取关联数据
	err := r.db.Preload("Product").
		Preload("Product.ProductImages").
		Preload("Buyer").
		Preload("Seller").
		First(&order, id).Error

	return &order, err
}

// CreateOrderLog 创建订单日志
func (r *OrderRepositoryImpl) CreateOrderLog(orderID uint, action, operator, remark string) error {
	log := models.OrderLog{
		OrderID:  orderID,
		Action:   action,
		Operator: operator,
		Remark:   remark,
	}

	return r.db.Create(&log).Error
}

// GetOrderLogs 获取订单日志
func (r *OrderRepositoryImpl) GetOrderLogs(orderID uint) ([]models.OrderLog, error) {
	var logs []models.OrderLog

	err := r.db.Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&logs).Error

	return logs, err
}
