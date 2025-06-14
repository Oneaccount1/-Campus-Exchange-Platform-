package repositories

import (
	"campus/internal/bootstrap"
	"github.com/casbin/casbin/v2"
)

// PermissionRepository 权限仓库接口
type PermissionRepository interface {
	// 角色管理
	AddRoleForUser(userID string, role string) error
	DeleteRoleForUser(userID string, role string) error
	GetRolesForUser(userID string) ([]string, error)
	HasRoleForUser(userID string, role string) (bool, error)

	// 权限管理
	AddPolicy(role string, obj string, act string) error
	RemovePolicy(role string, obj string, act string) error

	// 权限检查
	Enforce(userID string, obj string, act string) (bool, error)

	// 权限查询
	GetPermissionsForUser(userID string) ([][]string, error)
	GetPermissionsForRole(role string) ([][]string, error)

	// 批量操作
	AddPolicies(rules [][]string) error
	RemovePolicies(rules [][]string) error

	// 保存策略
	SavePolicy() error
}

// permissionRepository 权限仓库实现
type permissionRepository struct {
	enforcer *casbin.Enforcer
}

// NewPermissionRepository 创建权限仓库实例
func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{
		enforcer: bootstrap.GetEnforcer(),
	}
}

// AddRoleForUser 为用户添加角色
func (r *permissionRepository) AddRoleForUser(userID string, role string) error {
	_, err := r.enforcer.AddGroupingPolicy(userID, role)
	return err
}

// DeleteRoleForUser 移除用户角色
func (r *permissionRepository) DeleteRoleForUser(userID string, role string) error {
	_, err := r.enforcer.RemoveGroupingPolicy(userID, role)
	return err
}

// GetRolesForUser 获取用户角色
func (r *permissionRepository) GetRolesForUser(userID string) ([]string, error) {
	roles, err := r.enforcer.GetRolesForUser(userID)
	return roles, err
}

// HasRoleForUser 检查用户是否有特定角色
func (r *permissionRepository) HasRoleForUser(userID string, role string) (bool, error) {
	return r.enforcer.HasRoleForUser(userID, role)
}

// AddPolicy 添加权限策略
func (r *permissionRepository) AddPolicy(role string, obj string, act string) error {
	_, err := r.enforcer.AddPolicy(role, obj, act)
	return err
}

// RemovePolicy 移除权限策略
func (r *permissionRepository) RemovePolicy(role string, obj string, act string) error {
	_, err := r.enforcer.RemovePolicy(role, obj, act)
	return err
}

// Enforce 检查权限
func (r *permissionRepository) Enforce(userID string, obj string, act string) (bool, error) {
	return r.enforcer.Enforce(userID, obj, act)
}

// GetPermissionsForUser 获取用户权限
func (r *permissionRepository) GetPermissionsForUser(userID string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(userID)
}

// GetPermissionsForRole 获取角色权限
func (r *permissionRepository) GetPermissionsForRole(role string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(role)
}

// AddPolicies 批量添加策略
func (r *permissionRepository) AddPolicies(rules [][]string) error {
	_, err := r.enforcer.AddPolicies(rules)
	return err
}

// RemovePolicies 批量移除策略
func (r *permissionRepository) RemovePolicies(rules [][]string) error {
	_, err := r.enforcer.RemovePolicies(rules)
	return err
}

// SavePolicy 保存策略到存储
func (r *permissionRepository) SavePolicy() error {
	return r.enforcer.SavePolicy()
}
