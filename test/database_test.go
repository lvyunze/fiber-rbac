package test

import (
	"os"
	"testing"
	"time"

	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 在所有测试开始前执行
func TestMain(m *testing.M) {
	// 创建测试环境
	setup()

	// 运行测试
	code := m.Run()

	// 清理测试环境
	teardown()

	// 退出测试
	os.Exit(code)
}

// 初始化测试环境
func setup() {
	// 创建临时测试配置文件
	createTestConfig()
}

// 清理测试环境
func teardown() {
	// 删除测试使用的数据库文件
	os.Remove("test_db.db")
}

// 创建测试配置
func createTestConfig() {
	// 测试中使用内存数据库或临时文件
	os.Setenv("DATABASE_TYPE", "sqlite")
	os.Setenv("DATABASE_SQLITE_FILE", "test_db.db")
}

// 测试数据库初始化
func TestInitDB(t *testing.T) {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	err := database.InitDB(cfg)

	// 断言初始化不出错
	assert.NoError(t, err, "数据库初始化不应该出错")
	assert.NotNil(t, database.DB, "数据库连接不应该为空")

	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	assert.NoError(t, err, "获取sql.DB不应该出错")

	// 测试Ping
	err = sqlDB.Ping()
	assert.NoError(t, err, "数据库Ping测试不应该出错")

	// 测试数据库连接参数
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试获取连接参数
	maxIdle := sqlDB.Stats().Idle
	assert.LessOrEqual(t, maxIdle, 5, "最大空闲连接数应该小于等于5")
}

// 测试数据库事务
func TestDBTransaction(t *testing.T) {
	// 确保数据库已初始化
	if database.DB == nil {
		cfg := config.LoadConfig()
		err := database.InitDB(cfg)
		assert.NoError(t, err, "数据库初始化不应该出错")
	}

	// 测试事务
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些操作
		return nil
	})

	assert.NoError(t, err, "数据库事务不应该出错")
}

// 测试数据库关闭
func TestDBClose(t *testing.T) {
	// 确保数据库已初始化
	if database.DB == nil {
		cfg := config.LoadConfig()
		err := database.InitDB(cfg)
		assert.NoError(t, err, "数据库初始化不应该出错")
	}

	// 获取底层sql.DB连接
	sqlDB, err := database.DB.DB()
	assert.NoError(t, err, "获取sql.DB不应该出错")

	// 关闭数据库连接
	err = sqlDB.Close()
	assert.NoError(t, err, "关闭数据库连接不应该出错")
}
