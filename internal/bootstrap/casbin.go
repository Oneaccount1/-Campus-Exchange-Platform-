package bootstrap

import (
	"campus/internal/auth/casbin"
	"log"
)

// InitCasbin 初始化Casbin权限系统
func InitCasbin() error {
	// 初始化Casbin
	e, err := casbin.InitCasbin(GetDB())
	if err != nil {
		log.Printf("Casbin初始化失败: %v", err)
		return err
	}

	// 设置全局Casbin执行器
	SetEnforcer(e)

	log.Println("Casbin初始化成功")
	return nil
}
