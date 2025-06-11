package services

import (
	"campus/internal/bootstrap"
	"campus/internal/modules/permission/api"
	"campus/internal/utils/errors"
	"fmt"
)

// PermissionService 权限管理服务接口
type PermissionService interface {
	// 角色管理
	AddRoleForUser(userID uint, role string) error
	DeleteRoleForUser(userID uint, role string) error
	GetRolesForUser(userID uint) ([]string, error)

	// 权限管理
	AddPermissionForRole(role string, obj string, act string) error
	DeletePermissionForRole(role string, obj string, act string) error

	// 权限检查
	CheckPermission(userID uint, obj string, act string) (bool, error)

	// 获取用户权限列表
	GetUserPermissions(userID uint) (*api.PermissionListResponse, error)
}

type permissionService struct{}

// NewPermissionService 创建权限服务实例
func NewPermissionService() PermissionService {
	return &permissionService{}
}

// AddRoleForUser 为用户添加角色
func (s *permissionService) AddRoleForUser(userID uint, role string) error {
	enforcer := bootstrap.GetEnforcer()
	sub := fmt.Sprintf("%d", userID)
	_, err := enforcer.AddRoleForUser(sub, role)
	if err != nil {
		return errors.NewInternalServerError("添加角色失败", err)
	}
	return nil
}

// DeleteRoleForUser 删除用户的角色
func (s *permissionService) DeleteRoleForUser(userID uint, role string) error {
	enforcer := bootstrap.GetEnforcer()
	sub := fmt.Sprintf("%d", userID)
	_, err := enforcer.DeleteRoleForUser(sub, role)
	if err != nil {
		return errors.NewInternalServerError("删除角色失败", err)
	}
	return nil
}

// GetRolesForUser 获取用户的所有角色
func (s *permissionService) GetRolesForUser(userID uint) ([]string, error) {
	enforcer := bootstrap.GetEnforcer()
	sub := fmt.Sprintf("%d", userID)
	roles, err := enforcer.GetRolesForUser(sub)
	if err != nil {
		return nil, errors.NewInternalServerError("获取角色失败", err)
	}
	return roles, nil
}

// AddPermissionForRole 为角色添加权限
func (s *permissionService) AddPermissionForRole(role string, obj string, act string) error {
	enforcer := bootstrap.GetEnforcer()
	_, err := enforcer.AddPolicy(role, obj, act)
	if err != nil {
		return errors.NewInternalServerError("添加权限失败", err)
	}
	return nil
}

// DeletePermissionForRole 删除角色的权限
func (s *permissionService) DeletePermissionForRole(role string, obj string, act string) error {
	enforcer := bootstrap.GetEnforcer()
	_, err := enforcer.RemovePolicy(role, obj, act)
	if err != nil {
		return errors.NewInternalServerError("删除权限失败", err)
	}
	return nil
}

// CheckPermission 检查用户是否有权限
func (s *permissionService) CheckPermission(userID uint, obj string, act string) (bool, error) {
	enforcer := bootstrap.GetEnforcer()
	sub := fmt.Sprintf("%d", userID)
	return enforcer.Enforce(sub, obj, act)
}

// GetUserPermissions 获取用户的所有权限
func (s *permissionService) GetUserPermissions(userID uint) (*api.PermissionListResponse, error) {
	enforcer := bootstrap.GetEnforcer()
	sub := fmt.Sprintf("%d", userID)

	// 获取用户的所有角色
	roles, err := enforcer.GetRolesForUser(sub)
	if err != nil {
		return nil, errors.NewInternalServerError("获取角色失败", err)
	}

	// 用户的权限列表
	var permissions []api.Permission

	// 直接获取每个角色的权限，而不是获取所有策略再筛选
	for _, role := range roles {
		// 获取特定角色的权限策略
		rolePolicies, err := enforcer.GetFilteredPolicy(0, role)
		if err != nil {
			return nil, errors.NewInternalServerError("获取角色权限失败", err)
		}

		// 添加到权限列表
		for _, policy := range rolePolicies {
			if len(policy) >= 3 {
				permissions = append(permissions, api.Permission{
					Role:   role,
					Object: policy[1],
					Action: policy[2],
				})
			}
		}
	}

	// 获取用户直接分配的权限（不通过角色）
	userPolicies, err := enforcer.GetFilteredPolicy(0, sub)
	if err != nil {
		return nil, errors.NewInternalServerError("获取用户权限失败", err)
	}

	// 添加用户直接分配的权限
	for _, policy := range userPolicies {
		if len(policy) >= 3 {
			permissions = append(permissions, api.Permission{
				Role:   "direct", // 标记为直接分配的权限
				Object: policy[1],
				Action: policy[2],
			})
		}
	}

	return &api.PermissionListResponse{
		UserID:      userID,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}
