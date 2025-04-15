package model

import (
	"time"
	"gorm.io/gorm"
)

// UserRefreshToken 用户刷新令牌表结构
type UserRefreshToken struct {
	ID        uint64         `gorm:"primaryKey"`
	UserID    uint64         `gorm:"not null;index"`
	Token     string         `gorm:"size:512;not null;uniqueIndex"`
	ExpiresAt time.Time      `gorm:"not null;index"`
	Used      bool           `gorm:"default:false;not null"`
	Revoked   bool           `gorm:"default:false;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (UserRefreshToken) TableName() string {
	return "user_refresh_tokens"
}
