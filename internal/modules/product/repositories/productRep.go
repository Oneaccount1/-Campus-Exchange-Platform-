package repositories

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"campus/internal/modules/product/api"
	"campus/internal/utils/logger"
	"gorm.io/gorm"
	"time"
)

type ProductRepository interface {
	GetAll(page, size uint) ([]*models.Product, int64, error)
	GetByID(id string) (*models.Product, error)
	Create(product *models.Product) (uint, error)
	Update(id string, product *models.Product) error
	Delete(id string) error
	SearchProductsByKeyword(keyword string, page, size uint) ([]*models.Product, int64, error)
	GetByUserID(userID uint, page, size uint) ([]*models.Product, int64, error)
	GetSolvingProducts(page, size uint) ([]*models.Product, int64, error)
	BatchUpdateStatus(productIDs []uint, status string) error
	GetLatest(limit uint) ([]*models.Product, int64, error)
	FilterProducts(filter *api.FilterProductsRequest) ([]*models.Product, int64, error)
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func (r *ProductRepositoryImpl) FilterProducts(filter *api.FilterProductsRequest) ([]*models.Product, int64, error) {
	var products []*models.Product

	var tot int64

	// 构建查询
	query := r.db.Model(&models.Product{})
	// 过滤条件
	if filter.Keyword != "" {
		query = query.Where("title like ? OR description like ?", "%"+filter.Keyword, "%"+filter.Keyword)
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Condition != "" {
		query = query.Where("condition = ?", filter.Condition)
	}

	// 日期筛选
	if filter.StartDate != "" {
		date, err := time.Parse(time.RFC3339, filter.StartDate)
		if err != nil {
			logger.Errorf("开始日期格式化失败 : err %v", err)
		} else {
			query = query.Where("created_at >= ?", date)
		}
	}

	if filter.EndDate != "" {
		date, err := time.Parse(time.RFC3339, filter.EndDate)
		if err != nil {
			logger.Errorf("结束日期格式化失败 : err %v", err)
		} else {
			// 设置为当天最后一秒
			date = date.Add(24*time.Hour - time.Second)
			query = query.Where("created_at <= ?", date)
		}
	}

	if err := query.Count(&tot).Error; err != nil {
		return nil, 0, err
	}

	// 默认降序排序
	query = query.Order("created_at DESC")
	// 分页
	offset := (filter.Page - 1) * filter.Size
	query = query.Offset(int(offset)).Limit(int(filter.Size))
	query = query.Preload("ProductImages").Preload("User")
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, tot, nil
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{
		db: bootstrap.GetDB(),
	}
}

func (r *ProductRepositoryImpl) GetAll(page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	err := r.db.Model(&models.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	// 使用 Preload 预加载图片
	err = r.db.Preload("ProductImages").Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetSolvingProducts(page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	err := r.db.Model(&models.Product{}).Where("status = ?", "售卖中").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	// 使用 Preload 预加载图片
	err = r.db.Preload("ProductImages").Where("status = ?", "售卖中").Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetByID(id string) (*models.Product, error) {
	var product models.Product
	// 使用 Preload 预加载图片
	err := r.db.Preload("ProductImages").First(&product, "id = ?", id).Error
	return &product, err
}

func (r *ProductRepositoryImpl) Create(product *models.Product) (uint, error) {
	err := r.db.Create(product).Error
	return product.ID, err
}

func (r *ProductRepositoryImpl) Update(id string, product *models.Product) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(product).Error
}

func (r *ProductRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryImpl) SearchProductsByKeyword(keyword string, page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	query := "%" + keyword + "%"

	err := r.db.Model(&models.Product{}).Where("title LIKE ?", query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.Where("title LIKE ?", query).Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

func (r *ProductRepositoryImpl) GetByUserID(userID uint, page, size uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	// 计算总记录数
	err := r.db.Model(&models.Product{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err = r.db.Preload("ProductImages").Where("user_id = ?", userID).Offset(int(offset)).Limit(int(size)).Find(&products).Error
	return products, total, err
}

// BatchUpdateStatus 批量更新商品状态
func (r *ProductRepositoryImpl) BatchUpdateStatus(productIDs []uint, status string) error {
	return r.db.Model(&models.Product{}).Where("id IN ?", productIDs).Update("status", status).Error
}

// GetLatest 获取最新商品
func (r *ProductRepositoryImpl) GetLatest(limit uint) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	// 统计有效商品总数（状态为"售卖中"）
	err := r.db.Model(&models.Product{}).Where("status = ?", "售卖中").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取最新有效商品
	err = r.db.Preload("ProductImages").
		Preload("User").
		Where("status = ?", "售卖中").
		Order("created_at DESC").
		Limit(int(limit)).
		Find(&products).Error

	// 如果没有找到有效商品，尝试获取所有状态的最新商品
	//if len(products) == 0 {
	//	err = r.db.Preload("ProductImages").
	//		Preload("User").
	//		Order("created_at DESC").
	//		Limit(int(limit)).
	//		Find(&products).Error
	//}

	return products, total, err
}
