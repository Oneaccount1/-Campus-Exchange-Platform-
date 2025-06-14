package api

// RoleRequest 角色请求
type RoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// PermissionRequest 权限请求
type PermissionRequest struct {
	Role   string `json:"role" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// CheckPermissionRequest 检查权限请求
type CheckPermissionRequest struct {
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}
