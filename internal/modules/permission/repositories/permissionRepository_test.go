package repositories

import (
	"campus/internal/config"
	"campus/internal/database"
	"campus/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	casbinv2 "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// RealDatabaseTestSuite 使用真实数据库的测试套件
type RealDatabaseTestSuite struct {
	suite.Suite
	DB         *gorm.DB
	repository PermissionRepository
	enforcer   *casbinv2.Enforcer
	configPath string
	testRoles  []string
	testUsers  []string
	rootDir    string
	rbacPath   string
}

// SetupSuite 在所有测试之前运行
func (suite *RealDatabaseTestSuite) SetupSuite() {
	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		suite.T().Fatalf("获取工作目录失败: %v", err)
	}

	// 向上查找项目根目录
	rootDir := findRootDir(workDir)
	suite.rootDir = rootDir
	suite.configPath = filepath.Join(rootDir, "configs", "test_config.yaml")

	// 设置rbac.conf文件路径
	suite.rbacPath = filepath.Join(rootDir, "configs", "rbac_test.conf")

	// 创建测试用rbac.conf文件
	err = createDefaultRbacConf(suite.rbacPath)
	assert.NoError(suite.T(), err, "创建rbac.conf文件失败")

	// 加载测试配置
	cfg, err := config.LoadConfig(suite.configPath)
	assert.NoError(suite.T(), err, fmt.Sprintf("加载测试配置失败: %v", err))

	// 连接测试数据库
	db, err := database.NewDatabase(cfg.Database, cfg.Server.Mode)
	assert.NoError(suite.T(), err, fmt.Sprintf("连接测试数据库失败: %v", err))
	suite.DB = db

	// 迁移表结构
	err = suite.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	assert.NoError(suite.T(), err, fmt.Sprintf("迁移表结构失败: %v", err))

	// 初始化Casbin，直接自己初始化而不使用内部方法
	adapter, err := gormadapter.NewAdapterByDB(db)
	assert.NoError(suite.T(), err, fmt.Sprintf("创建GORM适配器失败: %v", err))

	enforcer, err := casbinv2.NewEnforcer(suite.rbacPath, adapter)
	assert.NoError(suite.T(), err, fmt.Sprintf("初始化Casbin Enforcer失败: %v", err))

	// 加载策略
	err = enforcer.LoadPolicy()
	assert.NoError(suite.T(), err, fmt.Sprintf("加载Casbin策略失败: %v", err))

	suite.enforcer = enforcer

	// 创建权限仓库实例
	suite.repository = &permissionRepository{
		enforcer: enforcer,
	}

	// 准备测试数据
	suite.prepareTestData()
}

// 创建默认的RBAC配置文件
func createDefaultRbacConf(path string) error {
	content := `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`

	return os.WriteFile(path, []byte(content), 0644)
}

// TearDownSuite 在所有测试之后运行
func (suite *RealDatabaseTestSuite) TearDownSuite() {
	// 清理测试数据
	for _, user := range suite.testUsers {
		suite.enforcer.DeleteUser(user)
	}
	for _, role := range suite.testRoles {
		suite.enforcer.DeleteRole(role)
	}

	// 删除所有与测试相关的策略
	suite.enforcer.RemoveFilteredPolicy(0, "test_role")
	suite.enforcer.RemoveFilteredPolicy(0, "test_admin_role")

	// 保存策略
	suite.enforcer.SavePolicy()

	// 删除测试配置文件
	if _, err := os.Stat(suite.rbacPath); err == nil {
		os.Remove(suite.rbacPath)
	}

	fmt.Println("测试套件清理完成")
}

// SetupTest 在每个测试之前运行
func (suite *RealDatabaseTestSuite) SetupTest() {
	// 清理上一次测试的数据
	for _, user := range suite.testUsers {
		suite.enforcer.DeleteUser(user)
	}
	for _, role := range suite.testRoles {
		suite.enforcer.DeleteRole(role)
	}

	// 清理所有测试策略
	suite.enforcer.RemoveFilteredPolicy(0, "test_role")
	suite.enforcer.RemoveFilteredPolicy(0, "test_admin_role")

	// 准备测试数据
	suite.prepareTestData()

	// 保存策略确保数据生效
	suite.enforcer.SavePolicy()
}

// TearDownTest 在每个测试之后运行
func (suite *RealDatabaseTestSuite) TearDownTest() {
	// 清理测试数据
	// 注意：这里不做完整清理，留给下一个测试的SetupTest处理
}

// 准备测试数据
func (suite *RealDatabaseTestSuite) prepareTestData() {
	// 设置测试角色
	suite.testRoles = []string{"test_role", "test_admin_role"}

	// 设置测试用户
	suite.testUsers = []string{"test_user1", "test_user2"}

	// 添加用户角色关系
	suite.enforcer.AddGroupingPolicy("test_user1", "test_role")
	suite.enforcer.AddGroupingPolicy("test_user2", "test_admin_role")

	// 添加角色权限策略
	suite.enforcer.AddPolicy("test_role", "/api/public", "GET")
	suite.enforcer.AddPolicy("test_admin_role", "/api/admin", "GET")
	suite.enforcer.AddPolicy("test_admin_role", "/api/admin", "POST")
}

// 查找项目根目录
func findRootDir(dir string) string {
	// 检查是否存在configs目录，如果存在则认为是项目根目录
	if _, err := os.Stat(filepath.Join(dir, "configs")); err == nil {
		return dir
	}

	// 向上一级目录查找
	parent := filepath.Dir(dir)
	if parent == dir {
		// 已经到达根目录，无法再向上
		return dir
	}
	return findRootDir(parent)
}

// TestAddRoleForUser 测试为用户添加角色
func (suite *RealDatabaseTestSuite) TestAddRoleForUser() {
	// 准备测试数据
	userID := "test_user_new"
	role := "test_role_new"

	// 测试添加角色
	err := suite.repository.AddRoleForUser(userID, role)
	assert.NoError(suite.T(), err, "添加用户角色失败")

	// 验证角色已添加
	hasRole, err := suite.repository.HasRoleForUser(userID, role)
	assert.NoError(suite.T(), err, "检查用户角色失败")
	assert.True(suite.T(), hasRole, "用户应该拥有刚添加的角色")

	// 清理测试数据
	suite.enforcer.DeleteUser(userID)
}

// TestDeleteRoleForUser 测试移除用户角色
func (suite *RealDatabaseTestSuite) TestDeleteRoleForUser() {
	// 准备测试数据
	userID := "test_delete_user"
	role := "test_delete_role"

	// 先添加角色，然后删除
	suite.enforcer.AddGroupingPolicy(userID, role)

	// 验证角色已添加
	hasRole, err := suite.repository.HasRoleForUser(userID, role)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasRole)

	// 测试删除角色
	err = suite.repository.DeleteRoleForUser(userID, role)
	assert.NoError(suite.T(), err)

	// 验证角色已删除
	hasRole, err = suite.repository.HasRoleForUser(userID, role)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasRole, "用户角色应该已被删除")

	// 清理
	suite.enforcer.DeleteUser(userID)
}

// TestGetRolesForUser 测试获取用户角色
func (suite *RealDatabaseTestSuite) TestGetRolesForUser() {
	// 已有测试数据
	userID := "test_user1"

	// 测试获取角色
	roles, err := suite.repository.GetRolesForUser(userID)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), roles, "test_role")
	assert.Len(suite.T(), roles, 1, "用户应该只有一个角色")

	// 测试不存在的用户
	roles, err = suite.repository.GetRolesForUser("nonexistent_user")
	assert.NoError(suite.T(), err, "对于不存在的用户，应返回空列表而不是错误")
	assert.Empty(suite.T(), roles, "不存在的用户应该没有角色")
}

// TestHasRoleForUser 测试检查用户角色
func (suite *RealDatabaseTestSuite) TestHasRoleForUser() {
	// 已有测试数据
	userID := "test_user2"
	role := "test_admin_role"

	// 测试拥有的角色
	hasRole, err := suite.repository.HasRoleForUser(userID, role)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasRole, "用户应该拥有该角色")

	// 测试不拥有的角色
	hasRole, err = suite.repository.HasRoleForUser(userID, "nonexistent_role")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasRole, "用户不应该拥有不存在的角色")
}

// TestAddPolicy 测试添加权限策略
func (suite *RealDatabaseTestSuite) TestAddPolicy() {
	// 准备测试数据
	role := "test_policy_role"
	obj := "/api/resource"
	act := "POST"

	// 测试添加策略
	err := suite.repository.AddPolicy(role, obj, act)
	assert.NoError(suite.T(), err)

	// 验证策略已添加
	policies, err := suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), policies, []string{role, obj, act})

	// 清理
	suite.enforcer.RemovePolicy(role, obj, act)
}

// TestRemovePolicy 测试移除权限策略
func (suite *RealDatabaseTestSuite) TestRemovePolicy() {
	// 准备测试数据
	role := "test_remove_policy_role"
	obj := "/api/test-resource"
	act := "GET"

	// 先添加策略
	suite.enforcer.AddPolicy(role, obj, act)

	// 验证策略已添加
	policies, err := suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), policies, []string{role, obj, act})

	// 测试移除策略
	err = suite.repository.RemovePolicy(role, obj, act)
	assert.NoError(suite.T(), err)

	// 验证策略已移除
	policies, err = suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	assert.NotContains(suite.T(), policies, []string{role, obj, act})
}

// TestEnforce 测试权限验证
func (suite *RealDatabaseTestSuite) TestEnforce() {
	// 已有测试数据
	userID := "test_user1"

	// 测试有权限的情况
	allowed, err := suite.repository.Enforce(userID, "/api/public", "GET")
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), allowed, "用户应该有权限访问公共API")

	// 测试无权限的情况
	allowed, err = suite.repository.Enforce(userID, "/api/admin", "GET")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), allowed, "普通用户不应该有权限访问管理员API")

	// 测试有权限的管理员
	allowed, err = suite.repository.Enforce("test_user2", "/api/admin", "GET")
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), allowed, "管理员应该有权限访问管理员API")
}

// TestGetPermissionsForUser 测试获取用户权限
func (suite *RealDatabaseTestSuite) TestGetPermissionsForUser() {
	// 已有测试数据
	userID := "test_user2"

	// 确保用户有角色，并且角色有权限
	roles, err := suite.repository.GetRolesForUser(userID)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), roles, "test_admin_role")

	// 验证角色有权限
	rolePolicies, err := suite.repository.GetPermissionsForRole("test_admin_role")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), rolePolicies)

	// 测试获取隐式权限（通过角色继承的权限）
	implicitPermissions, err := suite.enforcer.GetImplicitPermissionsForUser(userID)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), implicitPermissions, "管理员应该有继承的权限")

	// 验证特定的隐式权限
	foundGetPermission := false
	foundPostPermission := false
	for _, perm := range implicitPermissions {
		if len(perm) >= 3 && perm[1] == "/api/admin" && perm[2] == "GET" {
			foundGetPermission = true
		}
		if len(perm) >= 3 && perm[1] == "/api/admin" && perm[2] == "POST" {
			foundPostPermission = true
		}
	}
	assert.True(suite.T(), foundGetPermission, "管理员应该有GET /api/admin权限")
	assert.True(suite.T(), foundPostPermission, "管理员应该有POST /api/admin权限")
}

// TestGetPermissionsForRole 测试获取角色权限
func (suite *RealDatabaseTestSuite) TestGetPermissionsForRole() {
	// 已有测试数据
	role := "test_admin_role"

	// 测试获取权限
	permissions, err := suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), permissions, "管理员角色应该有权限")

	// 找到特定权限
	foundGetPermission := false
	foundPostPermission := false
	for _, perm := range permissions {
		if len(perm) >= 3 && perm[1] == "/api/admin" && perm[2] == "GET" {
			foundGetPermission = true
		}
		if len(perm) >= 3 && perm[1] == "/api/admin" && perm[2] == "POST" {
			foundPostPermission = true
		}
	}
	assert.True(suite.T(), foundGetPermission, "管理员角色应该有GET /api/admin权限")
	assert.True(suite.T(), foundPostPermission, "管理员角色应该有POST /api/admin权限")
}

// TestAddPolicies 测试批量添加策略
func (suite *RealDatabaseTestSuite) TestAddPolicies() {
	// 准备测试数据
	role := "test_batch_role"
	policies := [][]string{
		{role, "/api/batch1", "GET"},
		{role, "/api/batch2", "POST"},
		{role, "/api/batch3", "PUT"},
	}

	// 测试批量添加
	err := suite.repository.AddPolicies(policies)
	assert.NoError(suite.T(), err)

	// 验证所有策略都已添加
	rolePolicies, err := suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)

	for _, policy := range policies {
		assert.Contains(suite.T(), rolePolicies, policy, "应该包含添加的策略")
	}

	// 清理
	suite.enforcer.RemoveFilteredPolicy(0, role)
}

// TestRemovePolicies 测试批量移除策略
func (suite *RealDatabaseTestSuite) TestRemovePolicies() {
	// 准备测试数据
	role := "test_remove_batch_role"
	policies := [][]string{
		{role, "/api/remove1", "GET"},
		{role, "/api/remove2", "POST"},
	}

	// 先添加策略
	suite.enforcer.AddPolicies(policies)

	// 验证策略已添加
	rolePolicies, err := suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	for _, policy := range policies {
		assert.Contains(suite.T(), rolePolicies, policy)
	}

	// 测试批量移除
	err = suite.repository.RemovePolicies(policies)
	assert.NoError(suite.T(), err)

	// 验证策略已移除
	rolePolicies, err = suite.repository.GetPermissionsForRole(role)
	assert.NoError(suite.T(), err)
	for _, policy := range policies {
		assert.NotContains(suite.T(), rolePolicies, policy, "策略应该已被移除")
	}
}

// TestSavePolicy 测试保存策略
func (suite *RealDatabaseTestSuite) TestSavePolicy() {
	// 添加临时策略
	role := "test_save_role"
	obj := "/api/save-test"
	act := "GET"

	suite.enforcer.AddPolicy(role, obj, act)

	// 测试保存策略
	err := suite.repository.SavePolicy()
	assert.NoError(suite.T(), err)

	// 重新加载策略验证保存成功
	err = suite.enforcer.LoadPolicy()
	assert.NoError(suite.T(), err)

	// 验证策略仍然存在
	hasPolicy, err := suite.enforcer.HasPolicy(role, obj, act)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPolicy, "保存后策略应该仍然存在")

	// 清理
	suite.enforcer.RemovePolicy(role, obj, act)
	suite.enforcer.SavePolicy()
}

// TestRealDatabaseSuite 运行测试套件
func TestRealDatabaseSuite(t *testing.T) {
	suite.Run(t, new(RealDatabaseTestSuite))
}
