package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"size:32;not null;uniqueIndex" json:"username"`
	Email     string     `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string     `gorm:"size:255;not null" json:"-"` // 不输出到JSON
	CreatedAt int64      `gorm:"not null" json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
	DeletedAt *int64     `gorm:"index" json:"deleted_at"`
	Roles     []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 设置创建时间
	if u.CreatedAt == 0 {
		u.CreatedAt = NowUnix()
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 设置更新时间
	u.UpdatedAt = NowUnix()
	return nil
}

// UserRole 用户角色关联模型
type UserRole struct {
	UserID    uint64 `gorm:"primaryKey" json:"user_id"`
	RoleID    uint64 `gorm:"primaryKey" json:"role_id"`
	CreatedAt int64  `gorm:"not null;autoCreateTime" json:"created_at"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_roles"
}

// BeforeCreate 创建前钩子
func (ur *UserRole) BeforeCreate(tx *gorm.DB) error {
	if ur.CreatedAt == 0 {
		ur.CreatedAt = NowUnix()
	}
	return nil
}
