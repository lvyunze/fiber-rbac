package schema

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	Permission string `json:"permission" validate:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RoleIDs  []uint64 `json:"role_ids" validate:"omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       uint64   `json:"id" validate:"required"`
	Username string   `json:"username" validate:"required,min=3,max=32"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"omitempty,min=6"`
	RoleIDs  []uint64 `json:"role_ids" validate:"omitempty"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// GetUserRequest 获取用户详情请求
type GetUserRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// ListUserRequest 获取用户列表请求
type ListUserRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty"`
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint64        `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	CreatedAt int64         `json:"created_at"`
	Roles     []RoleSimple  `json:"roles,omitempty"`
}

// RoleSimple 简化的角色信息
type RoleSimple struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListUserResponse 用户列表响应
type ListUserResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}
