package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/database"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 准备基准测试数据库
func setupBenchmarkDB() (*gorm.DB, error) {
	// 使用内存数据库进行基准测试
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// 禁用日志以提高性能
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// 迁移数据库
	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 基准测试：数据库连接初始化
func BenchmarkInitDB(b *testing.B) {
	// 设置配置
	cfg := &config.Config{}
	cfg.Database.Type = "sqlite"
	cfg.Database.SQLite.File = "file::memory:?cache=shared"

	// 临时禁用标准日志输出
	oldOutput := log.Writer()
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(oldOutput)

	for i := 0; i < b.N; i++ {
		database.DB = nil
		_ = database.InitDB(cfg)
	}
}

// 基准测试：数据库查询
func BenchmarkDatabaseQueries(b *testing.B) {
	// 初始化测试数据库
	db, err := setupBenchmarkDB()
	if err != nil {
		b.Fatalf("设置测试数据库失败: %v", err)
	}

	// 创建测试用户
	testUser := &models.User{
		Username: "test_user",
		Password: "password123",
	}
	db.Create(testUser)

	// 基准测试查询
	b.Run("查询单个用户", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var user models.User
			db.First(&user, "username = ?", "test_user")
		}
	})

	// 基准测试创建记录
	b.Run("创建用户", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := models.User{
				Username: "user_" + time.Now().String(),
				Password: "password123",
			}
			db.Create(&user)
		}
	})

	// 基准测试更新记录
	b.Run("更新用户", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.Model(&models.User{}).Where("username = ?", "test_user").
				Update("username", "updated_user_"+time.Now().String())
		}
	})

	// 基准测试复杂查询
	b.Run("复杂查询", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var users []models.User
			db.Preload("Roles.Permissions").Find(&users)
		}
	})
}

// 基准测试：数据库事务
func BenchmarkDatabaseTransactions(b *testing.B) {
	// 初始化测试数据库
	db, err := setupBenchmarkDB()
	if err != nil {
		b.Fatalf("设置测试数据库失败: %v", err)
	}

	b.Run("简单事务", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.Transaction(func(tx *gorm.DB) error {
				user := models.User{
					Username: "tx_user_" + time.Now().String(),
					Password: "password123",
				}
				return tx.Create(&user).Error
			})
		}
	})

	b.Run("复杂事务", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.Transaction(func(tx *gorm.DB) error {
				// 创建用户
				user := models.User{
					Username: "complex_tx_user_" + time.Now().String(),
					Password: "password123",
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				// 创建角色
				role := models.Role{
					Name: "test_role_" + time.Now().String(),
				}
				if err := tx.Create(&role).Error; err != nil {
					return err
				}

				// 关联用户和角色
				return tx.Model(&user).Association("Roles").Append(&role)
			})
		}
	})
}

// 基准测试：连接池性能
func BenchmarkConnectionPool(b *testing.B) {
	// 初始化测试数据库
	db, err := setupBenchmarkDB()
	if err != nil {
		b.Fatalf("设置测试数据库失败: %v", err)
	}

	// 获取原始连接
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatalf("获取SQL DB失败: %v", err)
	}

	// 设置不同的连接池大小
	testPoolSizes := []int{5, 10, 20, 50}

	for _, size := range testPoolSizes {
		b.Run("连接池大小"+fmt.Sprintf("%d", size), func(b *testing.B) {
			sqlDB.SetMaxIdleConns(size)
			sqlDB.SetMaxOpenConns(size)

			// 模拟多个goroutine同时使用连接池
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					var user models.User
					db.First(&user, 1)
				}
			})
		})
	}
}
