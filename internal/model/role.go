package model

import (
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint64       `gorm:"primaryKey" json:"id"`
	Code        string       `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name        string       `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	CreatedAt   int64        `gorm:"not null" json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
	DeletedAt   *int64       `gorm:"index" json:"deleted_at"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate 创建前钩子
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	// 设置创建时间
	if r.CreatedAt == 0 {
		r.CreatedAt = NowUnix()
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	// 设置更新时间
	r.UpdatedAt = NowUnix()
	return nil
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	RoleID       uint64 `gorm:"primaryKey;not null" json:"role_id"`
	PermissionID uint64 `gorm:"primaryKey;not null" json:"permission_id"`
	CreatedAt    int64  `gorm:"not null" json:"created_at"`
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// BeforeCreate 创建前钩子
func (rp *RolePermission) BeforeCreate(tx *gorm.DB) error {
	// 设置创建时间
	if rp.CreatedAt == 0 {
		rp.CreatedAt = NowUnix()
	}
	return nil
}
