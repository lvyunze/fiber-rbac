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

// DeleteRoleRequest 
type DeleteRoleRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// GetRoleRequest 
type GetRoleRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// ListRoleRequest 
type ListRoleRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty"`
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

// ListRoleResponse 
type ListRoleResponse struct {
	Total int64          `json:"total"`
	Items []RoleResponse `json:"items"`
}
