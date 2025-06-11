package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// InitCasbin 初始化Casbin
func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer("configs/rbac.conf", adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	return e, nil
}
