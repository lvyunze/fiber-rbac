package schema

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID  uint64   `json:"user_id" validate:"required"`
	RoleIDs []uint64 `json:"role_ids" validate:"required"`
}

// GetUserRolesRequest 获取用户角色请求
type GetUserRolesRequest struct {
	UserID uint64 `json:"user_id" validate:"required"`
}

// UserRoleResponse 用户角色响应
type UserRoleResponse struct {
	UserID    uint64 `json:"user_id"`
	RoleID    uint64 `json:"role_id"`
	CreatedAt int64  `json:"created_at"`
}
