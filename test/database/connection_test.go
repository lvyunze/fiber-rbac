package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/database"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 数据库连接测试套件
type DBConnectionTestSuite struct {
	suite.Suite
	tempDBFile string
	config     *config.Config
}

// 手动加载测试配置
func loadTestConfig() (*config.Config, error) {
	viper.Reset()
	viper.SetConfigName("testconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./database")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading test config file: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

// 套件初始化
func (suite *DBConnectionTestSuite) SetupSuite() {
	// 创建临时数据库文件
	suite.tempDBFile = "test_db_suite.db"

	// 加载测试配置
	cfg, err := loadTestConfig()
	if err != nil {
		suite.T().Fatalf("无法加载测试配置: %v", err)
	}

	// 设置SQLite测试数据库文件
	cfg.Database.Type = "sqlite"
	cfg.Database.SQLite.File = suite.tempDBFile

	suite.config = cfg
}

// 套件清理
func (suite *DBConnectionTestSuite) TearDownSuite() {
	// 删除测试数据库文件
	os.Remove(suite.tempDBFile)
}

// 每个测试前的设置
func (suite *DBConnectionTestSuite) SetupTest() {
	// 确保每个测试前数据库连接是干净的
	database.DB = nil
}

// 测试数据库初始化
func (suite *DBConnectionTestSuite) TestInitDBWithSQLite() {
	// 使用SQLite配置
	suite.config.Database.Type = "sqlite"
	suite.config.Database.SQLite.File = suite.tempDBFile

	// 初始化数据库
	err := database.InitDB(suite.config)

	// 断言
	assert.NoError(suite.T(), err, "SQLite数据库初始化不应该出错")
	assert.NotNil(suite.T(), database.DB, "数据库连接不应该为空")

	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	assert.NoError(suite.T(), err, "获取sql.DB不应该出错")

	// 测试Ping
	err = sqlDB.Ping()
	assert.NoError(suite.T(), err, "数据库Ping测试不应该出错")

	// 关闭连接
	sqlDB.Close()
}

// 测试数据库迁移
func (suite *DBConnectionTestSuite) TestDatabaseMigration() {
	// 使用SQLite内存数据库进行测试
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(suite.T(), err, "创建内存数据库不应该出错")

	// 执行迁移
	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	assert.NoError(suite.T(), err, "数据库迁移不应该出错")

	// 测试表是否存在
	hasTable := db.Migrator().HasTable(&models.User{})
	assert.True(suite.T(), hasTable, "User表应该存在")

	hasTable = db.Migrator().HasTable(&models.Role{})
	assert.True(suite.T(), hasTable, "Role表应该存在")

	hasTable = db.Migrator().HasTable(&models.Permission{})
	assert.True(suite.T(), hasTable, "Permission表应该存在")
}

// 测试数据库连接池配置
func (suite *DBConnectionTestSuite) TestConnectionPoolConfig() {
	// 初始化数据库
	err := database.InitDB(suite.config)
	assert.NoError(suite.T(), err, "数据库初始化不应该出错")

	// 获取数据库连接
	sqlDB, err := database.DB.DB()
	assert.NoError(suite.T(), err, "获取sql.DB不应该出错")

	// 测试连接池配置
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 验证连接池配置
	maxIdle := sqlDB.Stats().Idle
	assert.LessOrEqual(suite.T(), maxIdle, 5, "最大空闲连接数应该小于等于5")

	// 关闭连接
	sqlDB.Close()
}

// 测试数据库错误处理
func (suite *DBConnectionTestSuite) TestDatabaseErrorHandling() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	assert.NoError(suite.T(), err, "获取工作目录不应该出错")

	// 设置错误的数据库配置
	invalidConfig := *suite.config
	invalidConfig.Database.Type = "sqlite"
	invalidConfig.Database.SQLite.File = filepath.Join(wd, "/non_existent_directory/test.db")

	// 尝试初始化数据库（应该返回错误）
	err = database.InitDB(&invalidConfig)
	assert.Error(suite.T(), err, "使用无效配置初始化数据库应该出错")

	// 测试错误消息
	if err != nil {
		assert.Contains(suite.T(), err.Error(), "failed to connect database", "错误消息应该包含'failed to connect database'")
	}
}

// 模拟不同类型数据库的测试
func (suite *DBConnectionTestSuite) TestDifferentDatabaseTypes() {
	// 测试SQLite
	testDB := func(dbType string) error {
		fmt.Printf("测试数据库类型: %s\n", dbType)
		testConfig := *suite.config
		testConfig.Database.Type = dbType

		// 根据数据库类型设置测试配置
		switch dbType {
		case "sqlite":
			testConfig.Database.SQLite.File = "file::memory:?cache=shared"
		case "mysql":
			// 这里不会真正连接MySQL，只是测试配置解析
			testConfig.Database.MySQL.Host = "localhost"
			testConfig.Database.MySQL.Port = 3306
			testConfig.Database.MySQL.User = "test_user"
			testConfig.Database.MySQL.Password = "test_password"
			testConfig.Database.MySQL.DBName = "test_db"
		case "postgres":
			// 这里不会真正连接PostgreSQL，只是测试配置解析
			testConfig.Database.Postgres.Host = "localhost"
			testConfig.Database.Postgres.Port = 5432
			testConfig.Database.Postgres.User = "test_user"
			testConfig.Database.Postgres.Password = "test_password"
			testConfig.Database.Postgres.DBName = "test_db"
			testConfig.Database.Postgres.SSLMode = "disable"
		}

		// 对于MySQL和PostgreSQL，我们只是测试配置解析，不实际连接
		if dbType == "mysql" || dbType == "postgres" {
			dsn := testConfig.GetDSN()
			assert.NotEmpty(suite.T(), dsn, fmt.Sprintf("%s DSN不应该为空", dbType))
			return nil
		}

		// 只对SQLite执行实际连接测试
		return database.InitDB(&testConfig)
	}

	// 测试SQLite连接
	err := testDB("sqlite")
	assert.NoError(suite.T(), err, "SQLite连接测试不应该出错")

	// 测试MySQL配置
	err = testDB("mysql")
	assert.NoError(suite.T(), err, "MySQL配置测试不应该出错")

	// 测试PostgreSQL配置
	err = testDB("postgres")
	assert.NoError(suite.T(), err, "PostgreSQL配置测试不应该出错")
}

// 运行测试套件
func TestDBConnectionSuite(t *testing.T) {
	suite.Run(t, new(DBConnectionTestSuite))
}
