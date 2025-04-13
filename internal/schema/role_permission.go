package schema

// AssignPermissionRequest 分配权限请求
type AssignPermissionRequest struct {
	RoleID       uint64   `json:"role_id" validate:"required"`
	PermissionIDs []uint64 `json:"permission_ids" validate:"required"`
}

// GetRolePermissionsRequest 获取角色权限请求
type GetRolePermissionsRequest struct {
	RoleID uint64 `json:"role_id" validate:"required"`
}

// RolePermissionResponse 角色权限响应
type RolePermissionResponse struct {
	RoleID       uint64 `json:"role_id"`
	PermissionID uint64 `json:"permission_id"`
	CreatedAt    int64  `json:"created_at"`
}
