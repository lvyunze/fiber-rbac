package model

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/lvyunze/fiber-rbac/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig, env string) error {
	// 根据环境设置数据库日志级别
	var logLevel logger.LogLevel
	var logLevelName string
	switch env {
	case "prod", "uat":
		// 生产和UAT环境只记录错误
		logLevel = logger.Error
		logLevelName = "ERROR"
	case "qa":
		// QA环境记录所有SQL
		logLevel = logger.Info
		logLevelName = "INFO"
	default:
		// 开发环境默认记录所有SQL
		logLevel = logger.Info
		logLevelName = "INFO"
	}

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logLevel), // 根据环境设置日志级别
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
	if err := db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
	); err != nil {
		return fmt.Errorf("自动迁移数据库模式失败: %w", err)
	}

	// 设置全局变量
	DB = db

	slog.Info("数据库连接初始化成功", "env", env, "log_level", logLevelName)
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
