package repositories

import (
	"campus/internal/database"
	"campus/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
	GetByID(id uint) (*models.Order, error)
	GetByBuyerID(buyerID uint, page, size uint) ([]*models.Order, int64, error)
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{
		db: database.GetDB(),
	}
}

func (r *OrderRepositoryImpl) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *OrderRepositoryImpl) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
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
