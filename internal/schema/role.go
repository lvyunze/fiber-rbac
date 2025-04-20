package schema

// CreateRoleRequest 
type CreateRoleRequest struct {
	Code        string   `json:"code" validate:"required,min=2,max=50"`
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	Description string   `json:"description" validate:"required"`
	PermissionIDs []uint64 `json:"permission_ids" validate:"omitempty"`
}

// UpdateRoleRequest 
type UpdateRoleRequest struct {
	ID          uint64   `json:"id" validate:"required"`
	Code        string   `json:"code" validate:"required,min=2,max=50"`
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	Description string   `json:"description" validate:"required"`
	PermissionIDs []uint64 `json:"permission_ids" validate:"omitempty"`
}

// RoleDeleteRequest 删除角色请求
type RoleDeleteRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// DeleteRoleRequest 
type DeleteRoleRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// GetRoleRequest 
type GetRoleRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// ListRoleRequest 角色列表请求参数
// 支持多字段排序、游标分页
// swagger:model
type ListRoleRequest struct {
	Page     int    `json:"page"`    // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Keyword  string `json:"keyword"`  // 关键词
	OrderBy  string `json:"order_by"` // 排序字段
	Desc     bool   `json:"desc"`     // 是否降序
}

// RoleDetailRequest 获取角色详情请求
type RoleDetailRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// RoleListPermissionsRequest 获取角色权限列表请求
type RoleListPermissionsRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// RoleResponse 
type RoleResponse struct {
	ID          uint64              `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CreatedAt   int64               `json:"created_at"`
	Permissions []PermissionSimple  `json:"permissions,omitempty"`
}

// PermissionSimple 
type PermissionSimple struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListRoleResponse 角色列表响应，包含分页信息
// Deprecated: use pagination.PageResult[RoleResponse]
type ListRoleResponse struct {
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	TotalPages int           `json:"total_pages"`
	Items     []RoleResponse `json:"items"`
}
