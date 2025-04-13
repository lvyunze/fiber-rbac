package model

import (
	"fmt"
	"log/slog"
	"github.com/lvyunze/fiber-rbac/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info), // 开发环境下设置详细日志
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 空闲连接最大生命周期

	// 自动迁移数据库模式
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("自动迁移数据库模式失败: %w", err)
	}

	// 设置全局DB实例
	DB = db
	slog.Info("数据库连接初始化成功")
	return nil
}

// 自动迁移数据库模式
func autoMigrate(db *gorm.DB) error {
	// 注册要迁移的模型
	return db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
	)
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
