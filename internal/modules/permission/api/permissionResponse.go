package api

// Permission 权限信息
type Permission struct {
	Role   string `json:"role"`
	Object string `json:"object"`
	Action string `json:"action"`
}

// PermissionListResponse 用户权限列表响应
type PermissionListResponse struct {
	UserID      uint         `json:"user_id"`
	Roles       []string     `json:"roles"`
	Permissions []Permission `json:"permissions"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Roles []string `json:"roles"`
}
