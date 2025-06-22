package services

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"campus/internal/modules/product/repositories"
	"campus/internal/modules/user/api"
	userRepo "campus/internal/modules/user/repositories"
	"campus/internal/utils/errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserService interface {
	Register(data *api.UserRegister) (*api.UserResponse, error)
	Login(data *api.UserLogin, roleCheck string) (*api.JWTResponse, error)
	GetByID(id uint) (*api.UserResponse, error)
	UpdateUser(id uint, dto api.UserUpdate) (*api.UserResponse, error)
	ChangePassword(id uint, oldPassword, newPassword string) error
	List(page, pageSize int) (*api.UserListResponse, error)
	AdminList(query *api.AdminUserListQuery) (*api.AdminUserListResponse, error)
	UpdateStatus(id uint, status string) error
	ResetPassword(id uint) (*api.ResetPasswordResponse, error)
	GetUserDetail(id uint) (*api.UserDetailResponse, error)
}
type userService struct {
	userRep    userRepo.UserRepository
	productRep repositories.ProductRepository
}

// convertToUserResponse 将User模型转换为UserResponse
func convertToUserResponse(user *models.User) *api.UserResponse {
	// 获取用户角色列表
	roleList := make([]string, 0)
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			roleList = append(roleList, role.Name)
		}
	}

	return &api.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Roles:        roleList,
		Description:  user.Description,
		Status:       user.Status,
		ProductCount: user.ProductCount,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// convertToAdminUserResponse 将User模型转换为AdminUserResponse
func convertToAdminUserResponse(user *models.User) *api.AdminUserResponse {
	return &api.AdminUserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar,
		RegisterTime: user.CreatedAt,
		ProductCount: user.ProductCount,
		Status:       user.Status,
		Email:        user.Email,
		Phone:        user.Phone,
	}
}

// Register 用户注册
func (u *userService) Register(data *api.UserRegister) (*api.UserResponse, error) {
	username, err := u.userRep.GetByUsername(data.Username)
	if err == nil && username != nil {
		return nil, errors.NewDuplicateError("用户名", nil)
	}
	email, err := u.userRep.GetByEmail(data.Email)
	if err == nil && email != nil {
		return nil, errors.NewDuplicateError("邮箱", nil)
	}
	// 加密密码
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败", err)
	}

	// 开启事务
	tx := bootstrap.GetDB().Begin()
	if tx.Error != nil {
		return nil, errors.NewInternalServerError("开启事务失败", tx.Error)
	}

	// 创建用户
	user := &models.User{
		Username:    data.Username,
		Password:    string(password),
		Email:       data.Email,
		Nickname:    data.Nickname,
		Phone:       data.Phone,
		Description: data.Description,
	}

	if err = tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalServerError("创建用户失败", err)
	}

	// 查找默认用户角色
	var userRole models.Role
	if err = tx.Where("name = ?", "user").First(&userRole).Error; err != nil {
		// 如果角色不存在，创建默认角色
		if err == gorm.ErrRecordNotFound {
			userRole = models.Role{
				Name:        "user",
				Description: "普通用户",
			}
			if err = tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalServerError("创建角色失败", err)
			}
		} else {
			tx.Rollback()
			return nil, errors.NewInternalServerError("查找角色失败", err)
		}
	}

	// 关联用户和角色
	if err = tx.Model(user).Association("Roles").Append(&userRole); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalServerError("关联角色失败", err)
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalServerError("提交事务失败", err)
	}

	// 加载用户角色关系
	if err = bootstrap.GetDB().Model(user).Association("Roles").Find(&user.Roles); err != nil {
		fmt.Printf("加载用户角色关系失败: %v\n", err)
	}

	return convertToUserResponse(user), nil
}

func (u *userService) Login(data *api.UserLogin, roleCheck string) (*api.JWTResponse, error) {
	user, err := u.userRep.GetByUsername(data.UserName)
	if err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误", err)
	}
	// 验证密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.PassWord)); err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误", err)
	}
	
	// 检查用户状态
	if user.Status == "禁用" {
		return nil, errors.NewForbiddenError("账号已被禁用，请联系管理员", nil)
	}

	// 获取JWT配置
	jwtConfig := bootstrap.GetConfig().JWT

	// 加载用户角色
	if err := bootstrap.GetDB().Model(user).Association("Roles").Find(&user.Roles); err != nil {
		// 记录错误但不影响登录流程
		fmt.Printf("加载用户角色关系失败: %v\n", err)
	}

	// 获取用户角色列表
	roleList := make([]string, 0)
	hasSpecificRole := roleCheck == "" // 如果不需要验证特定角色，默认为true

	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			roleList = append(roleList, role.Name)
			// 如果需要验证特定角色，检查用户是否拥有该角色
			if roleCheck != "" && role.Name == roleCheck {
				hasSpecificRole = true
			}
		}
	} else {
		// 如果没有关联角色，使用默认角色
		roleList = append(roleList, "user")
	}

	// 如果需要验证特定角色但用户没有该角色，返回未授权错误
	if !hasSpecificRole {
		return nil, errors.NewUnauthorizedError(fmt.Sprintf("您不是%s，无权访问", roleCheck), nil)
	}

	// 生成JWT
	expireTime := time.Now().Add(jwtConfig.Expiration)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roleList,
		"exp":      expireTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置中的密钥签名
	signedString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败", err)
	}

	return &api.JWTResponse{
		Token:     signedString,
		ExpiresAt: expireTime,
		UserID:    user.ID,
		Username:  user.Username,
		Roles:     roleList,
	}, nil
}

func (u *userService) GetByID(id uint) (*api.UserResponse, error) {
	user, err := u.userRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("用户", err)
	}

	// 加载用户角色
	if err := bootstrap.GetDB().Model(user).Association("Roles").Find(&user.Roles); err != nil {
		fmt.Printf("加载用户角色关系失败: %v\n", err)
	}

	// 获取用户产品数量
	_, total, err := u.productRep.GetByUserID(user.ID, 1, 1)
	if err != nil {
		fmt.Printf("获取用户产品数量失败: %v\n", err)
	} else {
		user.ProductCount = int(total)
	}

	return convertToUserResponse(user), nil
}

func (u *userService) UpdateUser(id uint, data api.UserUpdate) (*api.UserResponse, error) {
	user, err := u.userRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("用户", err)
	}
	// 检查邮箱是否被注册
	if data.Email != "" && data.Email != user.Email {
		existingUser, err := u.userRep.GetByEmail(data.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.NewDuplicateError("邮箱", nil)
		}
		user.Email = data.Email
	}
	// 修改其他信息
	if data.Nickname != "" {
		user.Nickname = data.Nickname
	}
	if data.Phone != "" {
		user.Phone = data.Phone
	}
	if data.Avatar != "" {
		user.Avatar = data.Avatar
	}
	if data.Description != "" {
		user.Description = data.Description
	}
	if err = u.userRep.Update(user); err != nil {
		return nil, errors.NewInternalServerError("更新用户信息失败", err)
	}

	// 加载用户角色
	if err := bootstrap.GetDB().Model(user).Association("Roles").Find(&user.Roles); err != nil {
		fmt.Printf("加载用户角色关系失败: %v\n", err)
	}

	return convertToUserResponse(user), nil
}

func (u *userService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := u.userRep.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("用户", err)
	}
	// 验证旧密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.NewBadRequestError("旧密码错误", err)
	}

	// 加密新密码
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError("密码加密失败", err)
	}
	user.Password = string(password)
	if err = u.userRep.Update(user); err != nil {
		return errors.NewInternalServerError("更新密码失败", err)
	}
	return nil
}

// List 获取用户列表
func (u *userService) List(page, pageSize int) (*api.UserListResponse, error) {
	users, total, err := u.userRep.List(page, pageSize)
	if err != nil {
		return nil, errors.NewInternalServerError("获取用户列表失败", err)
	}

	// 加载所有用户的角色
	for i := range users {
		if err := bootstrap.GetDB().Model(users[i]).Association("Roles").Find(&users[i].Roles); err != nil {
			fmt.Printf("加载用户角色关系失败: %v\n", err)
		}

		// 获取用户产品数量
		_, userProductTotal, err := u.productRep.GetByUserID(users[i].ID, 1, 1)
		if err != nil {
			fmt.Printf("获取用户产品数量失败: %v\n", err)
		} else {
			users[i].ProductCount = int(userProductTotal)
		}
	}

	userResponses := make([]api.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, *convertToUserResponse(user))
	}

	return &api.UserListResponse{
		List:  userResponses,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// AdminList 管理员获取用户列表（支持高级搜索）
func (u *userService) AdminList(query *api.AdminUserListQuery) (*api.AdminUserListResponse, error) {
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	// 调用仓库层的高级查询方法
	users, total, err := u.userRep.ListForAdmin(
		query.Page,
		query.Size,
		query.Search,
		query.Status,
		query.StartDate,
		query.EndDate,
	)

	if err != nil {
		return nil, errors.NewInternalServerError("获取用户列表失败", err)
	}

	// 处理每个用户的产品数量
	userResponses := make([]api.AdminUserResponse, 0, len(users))
	for _, user := range users {
		// 获取用户产品数量
		_, userProductTotal, err := u.productRep.GetByUserID(user.ID, 1, 1)
		if err != nil {
			fmt.Printf("获取用户产品数量失败: %v\n", err)
		} else {
			user.ProductCount = int(userProductTotal)
		}

		userResponses = append(userResponses, *convertToAdminUserResponse(user))
	}

	return &api.AdminUserListResponse{
		Total: total,
		List:  userResponses,
		Page:  query.Page,
		Size:  query.Size,
	}, nil
}

// UpdateStatus 更新用户状态
func (u *userService) UpdateStatus(id uint, status string) error {
	// 检查用户是否存在
	_, err := u.userRep.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("用户", err)
	}

	// 验证状态值
	if status != "正常" && status != "禁用" {
		return errors.NewBadRequestError("无效的状态值", nil)
	}

	// 更新状态
	return u.userRep.UpdateStatus(id, status)
}

// ResetPassword 重置用户密码
func (u *userService) ResetPassword(id uint) (*api.ResetPasswordResponse, error) {
	// 检查用户是否存在
	user, err := u.userRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("用户", err)
	}

	//// 生成6位随机密码
	//const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//const passwordLength = 6
	//
	//// 初始化随机数生成器
	//rand.Seed(time.Now().UnixNano())
	//
	//newPassword := make([]byte, passwordLength)
	//for i := range newPassword {
	//	newPassword[i] = charset[rand.Intn(len(charset))]
	//}
	newPassword := "123456"

	// 加密新密码
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败", err)
	}

	// 更新用户密码
	user.Password = string(password)
	if err = u.userRep.Update(user); err != nil {
		return nil, errors.NewInternalServerError("更新密码失败", err)
	}

	return &api.ResetPasswordResponse{
		NewPassword: string(newPassword),
	}, nil
}

// GetUserDetail 获取用户详情
func (u *userService) GetUserDetail(id uint) (*api.UserDetailResponse, error) {
	// 获取用户基本信息
	user, err := u.userRep.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("用户", err)
	}

	// 获取用户产品数量
	_, productTotal, err := u.productRep.GetByUserID(id, 1, 1)
	if err != nil {
		fmt.Printf("获取用户产品数量失败: %v\n", err)
		productTotal = 0
	}

	// 获取用户订单数量（买家）
	var orderCount int64
	if err := bootstrap.GetDB().Model(&models.Order{}).Where("buyer_id = ?", id).Count(&orderCount).Error; err != nil {
		fmt.Printf("获取用户订单数量失败: %v\n", err)
		orderCount = 0
	}

	// 获取用户收藏数量
	var favoriteCount int64
	if err := bootstrap.GetDB().Model(&models.Favorite{}).Where("user_id = ?", id).Count(&favoriteCount).Error; err != nil {
		fmt.Printf("获取用户收藏数量失败: %v\n", err)
		favoriteCount = 0
	}

	// 获取用户最近发布的商品
	products, _, err := u.productRep.GetByUserID(id, 1, 5)
	if err != nil {
		fmt.Printf("获取用户商品失败: %v\n", err)
		products = []*models.Product{}
	}

	userProducts := make([]api.UserProductItem, 0, len(products))
	for _, product := range products {
		userProducts = append(userProducts, api.UserProductItem{
			ID:         product.ID,
			Title:      product.Title,
			Price:      product.Price,
			CreateTime: product.CreatedAt,
		})
	}

	// 获取用户活动记录
	// 这里简单地组合用户的商品发布和订单记录
	activities := make([]api.UserActivityItem, 0)

	// 添加商品发布记录
	for _, product := range products {
		activities = append(activities, api.UserActivityItem{
			Content: fmt.Sprintf("发布了商品 \"%s\"", product.Title),
			Time:    product.CreatedAt,
		})
	}

	// 添加订单记录
	var orders []*models.Order
	if err := bootstrap.GetDB().Where("buyer_id = ? OR seller_id = ?", id, id).Order("created_at DESC").Limit(5).Find(&orders).Error; err == nil {
		for _, order := range orders {
			var product models.Product
			if err := bootstrap.GetDB().Select("title").First(&product, order.ProductID).Error; err == nil {
				var content string
				if order.BuyerID == id {
					content = fmt.Sprintf("购买了商品 \"%s\"", product.Title)
				} else {
					content = fmt.Sprintf("卖出了商品 \"%s\"", product.Title)
				}
				activities = append(activities, api.UserActivityItem{
					Content: content,
					Time:    order.CreatedAt,
				})
			}
		}
	}

	// 按时间排序活动（降序）
	// 简单冒泡排序
	for i := 0; i < len(activities)-1; i++ {
		for j := 0; j < len(activities)-i-1; j++ {
			if activities[j].Time.Before(activities[j+1].Time) {
				activities[j], activities[j+1] = activities[j+1], activities[j]
			}
		}
	}

	// 限制活动数量
	if len(activities) > 10 {
		activities = activities[:10]
	}

	return &api.UserDetailResponse{
		ID:            user.ID,
		Username:      user.Username,
		Avatar:        user.Avatar,
		RegisterTime:  user.CreatedAt,
		Email:         user.Email,
		Phone:         user.Phone,
		LastLogin:     user.UpdatedAt, // 暂时使用更新时间作为最后登录时间
		LastIP:        "",             // 暂无IP记录
		Status:        user.Status,
		ProductCount:  int(productTotal),
		OrderCount:    int(orderCount),
		FavoriteCount: int(favoriteCount),
		Products:      userProducts,
		Activities:    activities,
	}, nil
}

func NewUserService() UserService {
	return &userService{
		userRep:    userRepo.NewUserRepository(),
		productRep: repositories.NewProductRepository(),
	}
}
