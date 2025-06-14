package repositories

import (
	"campus/internal/config"
	"campus/internal/database"
	"campus/internal/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// RealDatabaseTestSuite 使用真实数据库的测试套件
type RealDatabaseTestSuite struct {
	suite.Suite
	DB         *gorm.DB
	repository UserRepository
	users      []*models.User
	roles      []*models.Role
	configPath string
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
	suite.configPath = filepath.Join(rootDir, "configs", "test_config.yaml")

	// 加载测试配置
	cfg, err := config.LoadConfig(suite.configPath)
	assert.NoError(suite.T(), err, fmt.Sprintf("加载测试配置失败: %v", err))

	// 连接测试数据库
	db, err := database.NewDatabase(cfg.Database, cfg.Server.Mode)
	assert.NoError(suite.T(), err, fmt.Sprintf("连接测试数据库失败: %v", err))
	suite.DB = db

	// 迁移表结构
	err = suite.DB.AutoMigrate(&models.User{}, &models.Role{})
	assert.NoError(suite.T(), err, fmt.Sprintf("迁移表结构失败: %v", err))

	// 创建仓库实例
	suite.repository = &userRepository{db: suite.DB}

	// 准备测试数据
	suite.prepareTestData()
}

// TearDownSuite 在所有测试之后运行
func (suite *RealDatabaseTestSuite) TearDownSuite() {
	// 清理测试数据
	suite.DB.Exec("DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'testuser%')")
	suite.DB.Exec("DELETE FROM users WHERE username LIKE 'testuser%'")
	suite.DB.Exec("DELETE FROM roles WHERE name LIKE 'test_%'")

	// 关闭数据库连接
	sqlDB, err := suite.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	fmt.Println("测试套件清理完成")
}

// SetupTest 在每个测试之前运行
func (suite *RealDatabaseTestSuite) SetupTest() {
	// 清空测试用户数据
	suite.DB.Exec("DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'testuser%')")
	suite.DB.Exec("DELETE FROM users WHERE username LIKE 'testuser%'")
	suite.DB.Exec("DELETE FROM roles WHERE name LIKE 'test_%'")

	// 创建角色
	for _, role := range suite.roles {
		err := suite.DB.Create(role).Error
		assert.NoError(suite.T(), err, fmt.Sprintf("插入角色测试数据失败"))
	}

	// 重新插入测试数据
	for _, user := range suite.users {
		err := suite.DB.Create(user).Error
		assert.NoError(suite.T(), err, fmt.Sprintf("插入用户测试数据失败"))
	}

	// 建立用户与角色的关联
	err := suite.DB.Exec("INSERT INTO user_roles (user_id, role_id) SELECT u.id, r.id FROM users u, roles r WHERE u.username = 'testuser1' AND r.name = 'test_user'").Error
	assert.NoError(suite.T(), err, "关联用户和角色失败")

	err = suite.DB.Exec("INSERT INTO user_roles (user_id, role_id) SELECT u.id, r.id FROM users u, roles r WHERE u.username = 'testuser2' AND r.name = 'test_admin'").Error
	assert.NoError(suite.T(), err, "关联用户和角色失败")
}

// TearDownTest 在每个测试之后运行
func (suite *RealDatabaseTestSuite) TearDownTest() {
	// 清空测试用户数据
	suite.DB.Exec("DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'testuser%')")
	suite.DB.Exec("DELETE FROM users WHERE username LIKE 'testuser%'")
	suite.DB.Exec("DELETE FROM roles WHERE name LIKE 'test_%'")
}

// 准备测试数据
func (suite *RealDatabaseTestSuite) prepareTestData() {
	// 创建测试角色
	suite.roles = []*models.Role{
		{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "test_user",
			Description: "Test user role",
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "test_admin",
			Description: "Test admin role",
		},
	}

	// 创建测试用户
	suite.users = []*models.User{
		{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username:    "testuser1",
			Password:    "password1",
			Nickname:    "Test User 1",
			Email:       "test1@example.com",
			Phone:       "1234567890",
			Avatar:      "avatar1.jpg",
			Description: "Test user 1 description",
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username:    "testuser2",
			Password:    "password2",
			Nickname:    "Test User 2",
			Email:       "test2@example.com",
			Phone:       "0987654321",
			Avatar:      "avatar2.jpg",
			Description: "Test user 2 description",
		},
	}
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

// 密码加密函数
func encryption(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
func equal(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// TestCreate 测试创建用户
func (suite *RealDatabaseTestSuite) TestCreate() {
	// 创建新用户
	newUser := &models.User{
		Username:    "testuser_new",
		Password:    "newpassword",
		Nickname:    "New Test User",
		Email:       "new_test@example.com",
		Phone:       "5555555555",
		Avatar:      "new-avatar.jpg",
		Description: "New test user description",
	}

	// 测试创建
	err := suite.repository.Create(newUser)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), newUser.ID, "用户ID应该被设置")

	// 验证用户已创建
	var foundUser models.User
	err = suite.DB.Where("username = ?", newUser.Username).First(&foundUser).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newUser.Username, foundUser.Username)
	bytes, err := encryption(newUser.Password)
	assert.NoError(suite.T(), err, fmt.Sprintf("加密失败！%v", err))
	assert.NoError(suite.T(), equal(string(bytes), foundUser.Password))

	assert.Equal(suite.T(), newUser.Email, foundUser.Email)
	assert.Equal(suite.T(), newUser.Nickname, foundUser.Nickname)
	assert.Equal(suite.T(), newUser.Phone, foundUser.Phone)
	assert.Equal(suite.T(), newUser.Description, foundUser.Description)

	// 关联用户角色
	userRole := suite.roles[0] // test_user角色
	err = suite.DB.Model(&newUser).Association("Roles").Append(userRole)
	assert.NoError(suite.T(), err, "关联用户角色失败")

	// 验证用户角色关联
	var userWithRoles models.User
	err = suite.DB.Preload("Roles").Where("username = ?", newUser.Username).First(&userWithRoles).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), userWithRoles.Roles, 1, "用户应该有一个角色")
	assert.Equal(suite.T(), "test_user", userWithRoles.Roles[0].Name)
}

// TestGetByID 测试通过ID获取用户
func (suite *RealDatabaseTestSuite) TestGetByID() {
	// 先获取第一个测试用户的ID
	var firstUser models.User
	err := suite.DB.Where("username = ?", "testuser1").First(&firstUser).Error
	assert.NoError(suite.T(), err)

	// 测试获取存在的用户
	user, err := suite.repository.GetByID(firstUser.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "testuser1", user.Username)

	// 测试获取不存在的用户
	user, err = suite.repository.GetByID(9999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

// TestGetByUsername 测试通过用户名获取用户
func (suite *RealDatabaseTestSuite) TestGetByUsername() {
	// 测试获取存在的用户
	user, err := suite.repository.GetByUsername("testuser1")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "testuser1", user.Username)

	// 测试获取不存在的用户
	user, err = suite.repository.GetByUsername("nonexistent")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

// TestGetByEmail 测试通过邮箱获取用户
func (suite *RealDatabaseTestSuite) TestGetByEmail() {
	// 测试获取存在的用户
	user, err := suite.repository.GetByEmail("test1@example.com")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "testuser1", user.Username)

	// 测试获取不存在的用户
	user, err = suite.repository.GetByEmail("nonexistent@example.com")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

// TestUpdate 测试更新用户
func (suite *RealDatabaseTestSuite) TestUpdate() {
	// 获取用户
	user, err := suite.repository.GetByUsername("testuser1")
	assert.NoError(suite.T(), err)

	// 更新用户信息
	user.Nickname = "Updated Nickname"
	user.Email = "updated@example.com"

	// 测试更新
	err = suite.repository.Update(user)
	assert.NoError(suite.T(), err)

	// 验证更新成功
	updatedUser, err := suite.repository.GetByUsername("testuser1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Nickname", updatedUser.Nickname)
	assert.Equal(suite.T(), "updated@example.com", updatedUser.Email)
}

// TestDelete 测试删除用户
func (suite *RealDatabaseTestSuite) TestDelete() {
	// 先获取第一个测试用户的ID
	var firstUser models.User
	err := suite.DB.Where("username = ?", "testuser1").First(&firstUser).Error
	assert.NoError(suite.T(), err)

	// 测试删除存在的用户
	err = suite.repository.Delete(firstUser.ID)
	assert.NoError(suite.T(), err)

	// 验证用户已删除
	_, err = suite.repository.GetByUsername("testuser1")
	assert.Error(suite.T(), err)

	// 测试删除不存在的用户
	err = suite.repository.Delete(9999)
	assert.NoError(suite.T(), err) // GORM在删除不存在的记录时不会返回错误
}

// TestList 测试列出用户
func (suite *RealDatabaseTestSuite) TestList() {
	// 测试第一页，每页1条记录
	users, total, err := suite.repository.List(1, 1)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), total, int64(2)) // 至少有2条测试记录
	assert.Equal(suite.T(), 1, len(users))            // 返回1条记录

	// 测试第二页，每页1条记录
	users, total, err = suite.repository.List(2, 1)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), total, int64(2)) // 至少有2条测试记录
	assert.Equal(suite.T(), 1, len(users))            // 返回1条记录

	// 测试每页10条记录
	users, total, err = suite.repository.List(1, 10)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), total, int64(2)) // 至少有2条测试记录
	assert.GreaterOrEqual(suite.T(), len(users), 2)   // 至少返回2条记录
}

// TestRealDatabaseSuite 运行测试套件
func TestRealDatabaseSuite(t *testing.T) {
	suite.Run(t, new(RealDatabaseTestSuite))
}

// TestInitRealDatabase 测试初始化真实数据库连接
func TestInitRealDatabase(t *testing.T) {
	// 获取工作目录
	workDir, err := os.Getwd()
	assert.NoError(t, err)

	// 向上查找项目根目录
	rootDir := findRootDir(workDir)
	configPath := filepath.Join(rootDir, "configs", "test_config.yaml")

	// 加载测试配置
	cfg, err := config.LoadConfig(configPath)
	assert.NoError(t, err)

	// 连接测试数据库
	db, err := database.NewDatabase(cfg.Database, cfg.Server.Mode)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// 关闭数据库连接
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	err = sqlDB.Close()
	assert.NoError(t, err)
}
