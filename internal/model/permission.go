package model

import (
	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	CreatedAt   int64  `gorm:"not null" json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   *int64 `gorm:"index" json:"deleted_at"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 设置表名
func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate 创建前钩子
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	// 设置创建时间
	if p.CreatedAt == 0 {
		p.CreatedAt = NowUnix()
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (p *Permission) BeforeUpdate(tx *gorm.DB) error {
	// 设置更新时间
	p.UpdatedAt = NowUnix()
	return nil
}
