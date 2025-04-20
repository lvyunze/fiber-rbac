package schema

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Code        string `json:"code" validate:"required,min=3,max=50"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	ID          uint64 `json:"id" validate:"required"`
	Code        string `json:"code" validate:"required,min=3,max=50"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required"`
}

// DeletePermissionRequest 删除权限请求
type DeletePermissionRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// GetPermissionRequest 获取权限详情请求
type GetPermissionRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// ListPermissionRequest 权限列表请求参数
// 支持多字段排序、游标分页
// swagger:model
type ListPermissionRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty"`
	OrderBy  string `json:"order_by"`
	Desc     bool   `json:"desc"`
}

// PermissionResponse 权限信息响应
type PermissionResponse struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

// ListPermissionResponse 权限列表响应，包含分页信息
// Deprecated: use pagination.PageResult[PermissionResponse]
type ListPermissionResponse struct {
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
	Items      []PermissionResponse `json:"items"`
}
