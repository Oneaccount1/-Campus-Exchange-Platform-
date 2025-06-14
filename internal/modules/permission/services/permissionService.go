package services

import (
	"campus/internal/modules/permission/api"
	"campus/internal/modules/permission/repositories"
	"campus/internal/utils/errors"
	"fmt"
	"strconv"
)

// PermissionService 权限服务接口
type PermissionService interface {
	// 角色管理
	AddRoleForUser(userID uint, role string) error
	DeleteRoleForUser(userID uint, role string) error
	GetRolesForUser(userID uint) ([]string, error)

	// 权限管理
	AddPermissionForRole(role, obj, act string) error
	DeletePermissionForRole(role, obj, act string) error

	// 权限检查
	CheckPermission(userID uint, obj, act string) (bool, error)

	// 获取用户权限列表
	GetUserPermissions(userID uint) (*api.PermissionListResponse, error)
}

// permissionService 权限服务实现
type permissionService struct {
	repo repositories.PermissionRepository
}

// NewPermissionService 创建权限服务实例
func NewPermissionService() PermissionService {
	return &permissionService{
		repo: repositories.NewPermissionRepository(),
	}
}

// AddRoleForUser 为用户添加角色
func (s *permissionService) AddRoleForUser(userID uint, role string) error {
	// 转换用户ID为字符串
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 使用仓库添加用户角色
	err := s.repo.AddRoleForUser(userIDStr, role)
	if err != nil {
		return errors.NewInternalServerError("添加角色失败", err)
	}

	// 保存策略到数据库
	err = s.repo.SavePolicy()
	if err != nil {
		return errors.NewInternalServerError("保存权限策略失败", err)
	}

	return nil
}

// DeleteRoleForUser 删除用户角色
func (s *permissionService) DeleteRoleForUser(userID uint, role string) error {
	// 转换用户ID为字符串
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 使用仓库移除用户角色
	err := s.repo.DeleteRoleForUser(userIDStr, role)
	if err != nil {
		return errors.NewInternalServerError("移除角色失败", err)
	}

	// 保存策略到数据库
	err = s.repo.SavePolicy()
	if err != nil {
		return errors.NewInternalServerError("保存权限策略失败", err)
	}

	return nil
}

// GetRolesForUser 获取用户角色
func (s *permissionService) GetRolesForUser(userID uint) ([]string, error) {
	// 转换用户ID为字符串
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 使用仓库获取用户角色
	roles, err := s.repo.GetRolesForUser(userIDStr)
	if err != nil {
		return nil, errors.NewInternalServerError("获取角色失败", err)
	}
	return roles, nil
}

// AddPermissionForRole 为角色添加权限
func (s *permissionService) AddPermissionForRole(role, obj, act string) error {
	// 使用仓库添加权限策略
	err := s.repo.AddPolicy(role, obj, act)
	if err != nil {
		return errors.NewInternalServerError("添加权限失败", err)
	}

	// 保存策略到数据库
	err = s.repo.SavePolicy()
	if err != nil {
		return errors.NewInternalServerError("保存权限策略失败", err)
	}

	return nil
}

// DeletePermissionForRole 删除角色权限
func (s *permissionService) DeletePermissionForRole(role, obj, act string) error {
	// 使用仓库移除权限策略
	err := s.repo.RemovePolicy(role, obj, act)
	if err != nil {
		return errors.NewInternalServerError("移除权限失败", err)
	}

	// 保存策略到数据库
	err = s.repo.SavePolicy()
	if err != nil {
		return errors.NewInternalServerError("保存权限策略失败", err)
	}

	return nil
}

// CheckPermission 检查用户是否有权限
func (s *permissionService) CheckPermission(userID uint, obj, act string) (bool, error) {
	// 转换用户ID为字符串
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 使用仓库检查权限
	hasPermission, err := s.repo.Enforce(userIDStr, obj, act)
	if err != nil {
		return false, errors.NewInternalServerError("权限检查失败", err)
	}
	return hasPermission, nil
}

// GetUserPermissions 获取用户所有权限
func (s *permissionService) GetUserPermissions(userID uint) (*api.PermissionListResponse, error) {
	// 获取用户角色
	roles, err := s.GetRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取所有权限策略
	allPermissions := make([]api.Permission, 0)

	// 遍历每个角色，获取其权限
	for _, role := range roles {
		permissions, err := s.repo.GetPermissionsForRole(role)
		if err != nil {
			return nil, errors.NewInternalServerError("获取角色权限失败", err)
		}

		// 转换为API响应格式
		for _, p := range permissions {
			if len(p) >= 3 {
				allPermissions = append(allPermissions, api.Permission{
					Role:   p[0],
					Object: p[1],
					Action: p[2],
				})
			}
		}
	}

	// 添加用户直接拥有的权限（如果有）
	userIDStr := fmt.Sprintf("%d", userID)
	userPermissions, err := s.repo.GetPermissionsForUser(userIDStr)
	if err != nil {
		return nil, errors.NewInternalServerError("获取用户权限失败", err)
	}

	for _, p := range userPermissions {
		if len(p) >= 3 {
			allPermissions = append(allPermissions, api.Permission{
				Role:   "direct",
				Object: p[1],
				Action: p[2],
			})
		}
	}

	return &api.PermissionListResponse{Permissions: allPermissions}, nil
}
