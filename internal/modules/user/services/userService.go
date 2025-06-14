package services

import (
	"campus/internal/bootstrap"
	"campus/internal/models"
	"campus/internal/modules/user/api"
	"campus/internal/modules/user/repositories"
	"campus/internal/utils/errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserService interface {
	Register(data *api.UserRegister) (*api.UserResponse, error)
	Login(data *api.UserLogin) (*api.JWTResponse, error)
	GetByID(id uint) (*api.UserResponse, error)
	UpdateUser(id uint, dto api.UserUpdate) (*api.UserResponse, error)
	ChangePassword(id uint, oldPassword, newPassword string) error
	List(page, pageSize int) (*api.UserListResponse, error)
}
type userService struct {
	userRep repositories.UserRepository
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
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		Roles:       roleList,
		Description: user.Description,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
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

func (u *userService) Login(data *api.UserLogin) (*api.JWTResponse, error) {
	user, err := u.userRep.GetByUsername(data.UserName)
	if err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误", err)
	}
	// 验证密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.PassWord)); err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误", err)
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
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			roleList = append(roleList, role.Name)
		}
	} else {
		// 如果没有关联角色，使用默认角色
		roleList = append(roleList, "user")
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
	}

	userResponses := make([]api.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, *convertToUserResponse(user))
	}

	return &api.UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

func NewUserService() UserService {
	return &userService{
		userRep: repositories.NewUserRepository(),
	}
}
