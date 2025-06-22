package services

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"campus/internal/modules/dashboard/api"
	"campus/internal/utils/errors"
	"fmt"
	"math"
	"time"
	"gorm.io/gorm"
)

// DashboardService 仪表盘服务接口
type DashboardService interface {
	// 获取仪表盘统计数据
	GetStats() (*api.DashboardStatsResponse, error)
	// 获取商品发布趋势
	GetProductTrend(timeRange string) (api.ProductTrendResponse, error)
	// 获取商品分类统计
	GetCategoryStats() (api.CategoryStatsResponse, error)
	// 获取最新商品
	GetLatestProducts(limit int) (api.LatestProductsResponse, error)
	// 获取系统活动
	GetActivities(limit int) (api.ActivitiesResponse, error)
}

// dashboardService 仪表盘服务实现
type dashboardService struct {
	db *gorm.DB
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService() DashboardService {
	return &dashboardService{
		db: bootstrap.GetDB(),
	}
}

// GetStats 获取仪表盘统计数据
func (s *dashboardService) GetStats() (*api.DashboardStatsResponse, error) {
	// 获取商品总数
	var productCount int64
	if err := s.db.Model(&models.Product{}).Count(&productCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取商品总数失败", err)
	}

	// 获取用户总数
	var userCount int64
	if err := s.db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取用户总数失败", err)
	}

	// 获取上周商品总数，计算趋势
	lastWeekTime := time.Now().AddDate(0, 0, -7)
	var lastWeekProductCount int64
	if err := s.db.Model(&models.Product{}).Where("created_at < ?", lastWeekTime).Count(&lastWeekProductCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取上周商品总数失败", err)
	}
	
	// 计算商品增长趋势
	var productTrend float64
	if lastWeekProductCount > 0 {
		productTrend = float64(productCount-lastWeekProductCount) / float64(lastWeekProductCount) * 100
	} else {
		productTrend = 100 // 如果上周没有商品，则增长率为100%
	}

	// 获取上周用户总数，计算趋势
	var lastWeekUserCount int64
	if err := s.db.Model(&models.User{}).Where("created_at < ?", lastWeekTime).Count(&lastWeekUserCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取上周用户总数失败", err)
	}
	
	// 计算用户增长趋势
	var userTrend float64
	if lastWeekUserCount > 0 {
		userTrend = float64(userCount-lastWeekUserCount) / float64(lastWeekUserCount) * 100
	} else {
		userTrend = 100 // 如果上周没有用户，则增长率为100%
	}

	// 本周新增用户数
	var newUserCount int64
	if err := s.db.Model(&models.User{}).Where("created_at >= ?", lastWeekTime).Count(&newUserCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取本周新增用户数失败", err)
	}

	// 今日新增商品数
	todayStart := time.Now().Truncate(24 * time.Hour)
	var todayProductCount int64
	if err := s.db.Model(&models.Product{}).Where("created_at >= ?", todayStart).Count(&todayProductCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取今日新增商品数失败", err)
	}

	// 昨日新增商品数
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	var yesterdayProductCount int64
	if err := s.db.Model(&models.Product{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Count(&yesterdayProductCount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取昨日新增商品数失败", err)
	}

	// 计算今日商品增长趋势
	var todayProductTrend float64
	if yesterdayProductCount > 0 {
		todayProductTrend = float64(todayProductCount-yesterdayProductCount) / float64(yesterdayProductCount) * 100
	} else if todayProductCount > 0 {
		todayProductTrend = 100 // 如果昨天没有商品而今天有，则增长率为100%
	} else {
		todayProductTrend = 0 // 如果昨天和今天都没有商品，则增长率为0
	}

	// 交易总额
	var totalAmount float64
	if err := s.db.Model(&models.Order{}).
		Joins("JOIN products ON orders.product_id = products.id").
		Where("orders.status = ?", "卖家已同意").
		Select("COALESCE(SUM(products.price), 0)").
		Scan(&totalAmount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取交易总额失败", err)
	}

	// 上月交易总额
	lastMonthStart := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	lastMonthEnd := time.Now().Truncate(24 * time.Hour)
	var lastMonthAmount float64
	if err := s.db.Model(&models.Order{}).
		Joins("JOIN products ON orders.product_id = products.id").
		Where("orders.status = ? AND orders.created_at >= ? AND orders.created_at < ?", "卖家已同意", lastMonthStart, lastMonthEnd).
		Select("COALESCE(SUM(products.price), 0)").
		Scan(&lastMonthAmount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取上月交易总额失败", err)
	}

	// 计算交易额增长趋势
	var amountTrend float64
	if lastMonthAmount > 0 {
		amountTrend = (totalAmount - lastMonthAmount) / lastMonthAmount * 100
	} else if totalAmount > 0 {
		amountTrend = 100 // 如果上月没有交易而本月有，则增长率为100%
	} else {
		amountTrend = 0 // 如果上月和本月都没有交易，则增长率为0
	}

	// 本月交易额
	thisMonthStart := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	var monthAmount float64
	if err := s.db.Model(&models.Order{}).
		Joins("JOIN products ON orders.product_id = products.id").
		Where("orders.status = ? AND orders.created_at >= ?", "卖家已同意", thisMonthStart).
		Select("COALESCE(SUM(products.price), 0)").
		Scan(&monthAmount).Error; err != nil {
		return nil, errors.NewInternalServerError("获取本月交易额失败", err)
	}

	// 返回结果
	return &api.DashboardStatsResponse{
		ProductCount:          int(productCount),
		ProductTrend:          math.Round(productTrend*100) / 100, // 保留两位小数
		UserCount:             int(userCount),
		UserTrend:             math.Round(userTrend*100) / 100,
		NewUserCount:          int(newUserCount),
		TodayProductCount:     int(todayProductCount),
		TodayProductTrend:     math.Round(todayProductTrend*100) / 100,
		YesterdayProductCount: int(yesterdayProductCount),
		TotalAmount:           math.Round(totalAmount*100) / 100,
		AmountTrend:           math.Round(amountTrend*100) / 100,
		MonthAmount:           math.Round(monthAmount*100) / 100,
	}, nil
}

// GetProductTrend 获取商品发布趋势
func (s *dashboardService) GetProductTrend(timeRange string) (api.ProductTrendResponse, error) {
	var result api.ProductTrendResponse

	// 根据时间范围设置不同的查询
	var startTime time.Time
	var labels []string
	
	if timeRange == "month" {
		// 过去30天的数据，按天统计
		startTime = time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
		for i := 0; i < 30; i++ {
			date := startTime.AddDate(0, 0, i)
			labels = append(labels, date.Format("01-02"))
		}
	} else {
		// 默认为week，过去7天的数据，按天统计
		startTime = time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
		for i := 0; i < 7; i++ {
			date := startTime.AddDate(0, 0, i)
			labels = append(labels, weekdays[date.Weekday()])
		}
	}

	for i, label := range labels {
		date := startTime.AddDate(0, 0, i)
		nextDate := date.AddDate(0, 0, 1)

		// 查询当天创建的商品数量
		var count int64
		if err := s.db.Model(&models.Product{}).Where("created_at >= ? AND created_at < ?", date, nextDate).Count(&count).Error; err != nil {
			return nil, errors.NewInternalServerError(fmt.Sprintf("获取%s商品数量失败", label), err)
		}

		result = append(result, api.ProductTrendItem{
			Label: label,
			Value: int(count),
		})
	}

	return result, nil
}

// GetCategoryStats 获取商品分类统计
func (s *dashboardService) GetCategoryStats() (api.CategoryStatsResponse, error) {
	var result api.CategoryStatsResponse

	// 查询所有分类及其商品数量
	rows, err := s.db.Model(&models.Product{}).
		Select("category, COUNT(*) as count").
		Where("category != ''").
		Group("category").
		Order("count DESC").
		Rows()
	
	if err != nil {
		return nil, errors.NewInternalServerError("获取分类统计失败", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, errors.NewInternalServerError("解析分类统计失败", err)
		}
		
		result = append(result, api.CategoryStatsItem{
			Name:  category,
			Value: count,
		})
	}

	return result, nil
}

// GetLatestProducts 获取最新商品
func (s *dashboardService) GetLatestProducts(limit int) (api.LatestProductsResponse, error) {
	var result api.LatestProductsResponse

	// 如果没有指定limit，默认为5
	if limit <= 0 {
		limit = 5
	}

	var products []*models.Product
	if err := s.db.Preload("ProductImages").Order("created_at DESC").Limit(limit).Find(&products).Error; err != nil {
		return nil, errors.NewInternalServerError("获取最新商品失败", err)
	}

	for _, product := range products {
		// 获取卖家信息
		var seller models.User
		if err := s.db.Select("username").First(&seller, product.UserID).Error; err != nil {
			// 如果找不到卖家，使用"未知用户"
			seller.Username = "未知用户"
		}

		// 获取商品图片
		var imageURL string
		if len(product.ProductImages) > 0 {
			imageURL = product.ProductImages[0].ImageURL
		}

		result = append(result, api.LatestProductItem{
			ID:         product.ID,
			Title:      product.Title,
			Price:      product.Price,
			Category:   product.Category,
			Seller:     seller.Username,
			CreateTime: product.CreatedAt,
			Image:      imageURL,
			Status:     product.Status,
		})
	}

	return result, nil
}

// GetActivities 获取系统活动
func (s *dashboardService) GetActivities(limit int) (api.ActivitiesResponse, error) {
	var result api.ActivitiesResponse

	// 如果没有指定limit，默认为5
	if limit <= 0 {
		limit = 5
	}

	// 查询最近的商品创建记录
	var products []*models.Product
	if err := s.db.Order("created_at DESC").Limit(limit).Find(&products).Error; err != nil {
		return nil, errors.NewInternalServerError("获取最近商品活动失败", err)
	}

	// 查询最近的订单记录
	var orders []*models.Order
	if err := s.db.Order("created_at DESC").Limit(limit).Find(&orders).Error; err != nil {
		return nil, errors.NewInternalServerError("获取最近订单活动失败", err)
	}

	// 组合商品和订单活动
	for _, product := range products {
		// 获取卖家信息
		var seller models.User
		if err := s.db.Select("username").First(&seller, product.UserID).Error; err != nil {
			// 如果找不到卖家，使用"未知用户"
			seller.Username = "未知用户"
		}

		timeAgo := getTimeAgo(product.CreatedAt)
		result = append(result, api.ActivityItem{
			Content: fmt.Sprintf("用户 %s 发布了商品 \"%s\"", seller.Username, product.Title),
			Time:    timeAgo,
			Type:    "primary",
			Color:   "#409EFF",
		})
	}

	for _, order := range orders {
		// 获取买家和卖家信息
		var buyer, seller models.User
		if err := s.db.Select("username").First(&buyer, order.BuyerID).Error; err != nil {
			buyer.Username = "未知用户"
		}
		if err := s.db.Select("username").First(&seller, order.SellerID).Error; err != nil {
			seller.Username = "未知用户"
		}

		// 获取商品信息
		var product models.Product
		if err := s.db.Select("title").First(&product, order.ProductID).Error; err != nil {
			product.Title = "未知商品"
		}

		timeAgo := getTimeAgo(order.CreatedAt)
		color := "#E6A23C"
		typeStr := "warning"
		
		// 根据订单状态设置不同的颜色和类型
		if order.Status == "卖家已同意" {
			color = "#67C23A"
			typeStr = "success"
		} else if order.Status == "卖家已拒绝" {
			color = "#F56C6C"
			typeStr = "danger"
		}

		result = append(result, api.ActivityItem{
			Content: fmt.Sprintf("用户 %s 购买了 %s 的商品 \"%s\"", buyer.Username, seller.Username, product.Title),
			Time:    timeAgo,
			Type:    typeStr,
			Color:   color,
		})
	}

	// 按时间排序并限制数量
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// getTimeAgo 获取时间差的友好描述
func getTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "刚刚"
	case diff < time.Hour:
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	case diff < 48*time.Hour:
		return "昨天"
	case diff < 72*time.Hour:
		return "前天"
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d天前", int(diff.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
} 