package bootstrap

import (
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

		log.Println("默认权限初始化成功")
	} else {
		log.Println("跳过权限初始化，使用数据库中现有权限")
	}

	return nil
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
