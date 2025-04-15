package test

import (
	"testing"
	"time"
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lvyunze/fiber-rbac/internal/model"
)

// TestUserRefreshTokenModelPG 测试 user_refresh_tokens 表的基本 CRUD（PostgreSQL，使用项目 config）
func TestUserRefreshTokenModelPG(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	assert.NoError(t, err)
	dsn := cfg.Database.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// 清理表
	db.Exec("DROP TABLE IF EXISTS user_refresh_tokens")

	// 自动迁移
	err = db.AutoMigrate(&model.UserRefreshToken{})
	assert.NoError(t, err)

	// 新增
	rt := &model.UserRefreshToken{
		UserID:    1,
		Token:     "test-token",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	assert.NoError(t, db.Create(rt).Error)

	// 查询
	var got model.UserRefreshToken
	err = db.Where("token = ?", "test-token").First(&got).Error
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), got.UserID)
	assert.False(t, got.Used)
	assert.False(t, got.Revoked)

	// 更新 used
	assert.NoError(t, db.Model(&model.UserRefreshToken{}).Where("token = ?", "test-token").Update("used", true).Error)
	var updated model.UserRefreshToken
	_ = db.Where("token = ?", "test-token").First(&updated).Error
	assert.True(t, updated.Used)

	// 撤销
	assert.NoError(t, db.Model(&model.UserRefreshToken{}).Where("token = ?", "test-token").Update("revoked", true).Error)
	_ = db.Where("token = ?", "test-token").First(&updated).Error
	assert.True(t, updated.Revoked)
}
