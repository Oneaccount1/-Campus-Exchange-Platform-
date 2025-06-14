package bootstrap

import (
	"campus/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

// InitPermissions 初始化权限
func InitPermissions() error {
	// 获取Enforcer
	enforcer := GetEnforcer()
	if enforcer == nil {
		log.Println("Enforcer未初始化，无法设置默认权限")
		return nil
	}

	// 获取配置
	config := GetConfig()

	// 只在特定条件下初始化权限：
	// 1. 强制初始化权限选项为true
	// 2. 或者是debug模式且数据库中没有权限规则
	if config.Server.InitPermissions || (config.Server.Mode == "debug" && shouldInitPermissions()) {
		log.Println("初始化默认权限...")

		// 清除所有现有策略
		enforcer.ClearPolicy()

		// 添加基本权限
		// 用户权限
		enforcer.AddPolicy("user", "/api/v1/user/profile", "GET")          // 获取个人资料
		enforcer.AddPolicy("user", "/api/v1/user/profile", "PUT")          // 更新个人资料
		enforcer.AddPolicy("user", "/api/v1/user/change-password", "POST") // 修改密码
		enforcer.AddPolicy("user", "/api/v1/user/:id", "GET")              // 查看用户信息

		// 管理员权限
		enforcer.AddPolicy("admin", "/api/v1/admin/users", "GET")                             // 获取用户列表
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/users/:id/roles", "POST")      // 分配角色
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/users/:id/roles", "DELETE")    // 移除角色
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/users/:id/roles", "GET")       // 获取用户角色
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/policies", "POST")             // 添加权限
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/policies", "DELETE")           // 移除权限
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/check", "POST")                // 检查权限
		enforcer.AddPolicy("admin", "/api/v1/admin/permissions/users/:id/permissions", "GET") // 获取用户权限

		// 保存策略到数据库
		err := enforcer.SavePolicy()
		if err != nil {
			log.Printf("保存权限策略失败: %v", err)
			return err
		}

		// 初始化数据库中的角色和权限
		initRolesAndPermissions()

		log.Println("默认权限初始化成功")
	} else {
		log.Println("跳过权限初始化，使用数据库中现有权限")
	}

	return nil
}

// initRolesAndPermissions 初始化数据库中的角色和权限
func initRolesAndPermissions() {
	db := GetDB()
	if db == nil {
		log.Println("数据库未初始化，无法创建默认角色和权限")
		return
	}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("开启事务失败: %v", tx.Error)
		return
	}

	// 创建默认角色
	roles := []models.Role{
		{Name: "user", Description: "普通用户"},
		{Name: "admin", Description: "管理员"},
	}

	for _, role := range roles {
		var count int64
		tx.Model(&models.Role{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			if err := tx.Create(&role).Error; err != nil {
				tx.Rollback()
				log.Printf("创建角色失败: %v", err)
				return
			}
		}
	}

	// 创建默认权限
	permissions := []models.Permission{
		{Resource: "/api/v1/user/profile", Action: "GET", Description: "获取个人资料"},
		{Resource: "/api/v1/user/profile", Action: "PUT", Description: "更新个人资料"},
		{Resource: "/api/v1/user/change-password", Action: "POST", Description: "修改密码"},
		{Resource: "/api/v1/user/:id", Action: "GET", Description: "查看用户信息"},
		{Resource: "/api/v1/admin/users", Action: "GET", Description: "获取用户列表"},
		{Resource: "/api/v1/admin/permissions/users/:id/roles", Action: "POST", Description: "分配角色"},
		{Resource: "/api/v1/admin/permissions/users/:id/roles", Action: "DELETE", Description: "移除角色"},
		{Resource: "/api/v1/admin/permissions/users/:id/roles", Action: "GET", Description: "获取用户角色"},
		{Resource: "/api/v1/admin/permissions/policies", Action: "POST", Description: "添加权限"},
		{Resource: "/api/v1/admin/permissions/policies", Action: "DELETE", Description: "移除权限"},
		{Resource: "/api/v1/admin/permissions/check", Action: "POST", Description: "检查权限"},
		{Resource: "/api/v1/admin/permissions/users/:id/permissions", Action: "GET", Description: "获取用户权限"},
	}

	for _, perm := range permissions {
		var count int64
		tx.Model(&models.Permission{}).Where("resource = ? AND action = ?", perm.Resource, perm.Action).Count(&count)
		if count == 0 {
			if err := tx.Create(&perm).Error; err != nil {
				tx.Rollback()
				log.Printf("创建权限失败: %v", err)
				return
			}
		}
	}

	// 关联角色和权限
	var userRole models.Role
	var adminRole models.Role

	tx.Where("name = ?", "user").First(&userRole)
	tx.Where("name = ?", "admin").First(&adminRole)

	// 用户权限
	var userPermissions []models.Permission
	tx.Where("resource IN (?, ?, ?, ?)",
		"/api/v1/user/profile",
		"/api/v1/user/profile",
		"/api/v1/user/change-password",
		"/api/v1/user/:id").Find(&userPermissions)

	if len(userPermissions) > 0 {
		if err := tx.Model(&userRole).Association("Permissions").Append(&userPermissions); err != nil {
			tx.Rollback()
			log.Printf("关联用户角色权限失败: %v", err)
			return
		}
	}

	// 管理员权限（包括用户所有权限）
	var adminPermissions []models.Permission
	tx.Find(&adminPermissions) // 所有权限

	if len(adminPermissions) > 0 {
		if err := tx.Model(&adminRole).Association("Permissions").Append(&adminPermissions); err != nil {
			tx.Rollback()
			log.Printf("关联管理员角色权限失败: %v", err)
			return
		}
	}

	// 检查是否有admin用户，如果没有，创建一个默认管理员账户
	var adminCount int64
	tx.Model(&models.User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		// 创建默认管理员
		createDefaultAdmin(tx)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("提交事务失败: %v", err)
		return
	}

	log.Println("默认角色和权限初始化成功")
}

// createDefaultAdmin 创建默认管理员账户
func createDefaultAdmin(tx *gorm.DB) {
	// 加密默认密码
	password, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("加密管理员密码失败: %v", err)
		return
	}

	// 创建管理员用户
	admin := &models.User{
		Username:    "admin",
		Password:    string(password),
		Email:       "admin@example.com",
		Nickname:    "系统管理员",
		Description: "系统默认管理员账户",
	}

	if err := tx.Create(admin).Error; err != nil {
		log.Printf("创建管理员账户失败: %v", err)
		return
	}

	// 关联管理员角色
	var adminRole models.Role
	if err := tx.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("查找管理员角色失败: %v", err)
		return
	}

	if err := tx.Model(admin).Association("Roles").Append(&adminRole); err != nil {
		log.Printf("关联管理员角色失败: %v", err)
		return
	}

	log.Println("创建默认管理员账户成功")
}

// shouldInitPermissions 判断是否应该初始化权限
// 这里可以实现检查数据库中是否已有权限规则的逻辑
func shouldInitPermissions() bool {
	// 获取Enforcer
	enforcer := GetEnforcer()

	// 检查是否有任何策略
	policies, _ := enforcer.GetPolicy()

	// 如果没有策略，则需要初始化
	return len(policies) == 0
}
