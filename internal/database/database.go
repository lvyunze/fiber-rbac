package database

import (
	"fmt"
	"log"

	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(getDialector(cfg), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func getDialector(cfg *config.Config) gorm.Dialector {
	switch cfg.Database.Type {
	case "sqlite":
		return sqlite.Open(cfg.Database.SQLite.File)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.MySQL.User,
			cfg.Database.MySQL.Password,
			cfg.Database.MySQL.Host,
			cfg.Database.MySQL.Port,
			cfg.Database.MySQL.DBName,
		)
		return mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.DBName,
			cfg.Database.Postgres.SSLMode,
		)
		return postgres.Open(dsn)
	default:
		log.Printf("Unknown database type: %s, using SQLite as default", cfg.Database.Type)
		return sqlite.Open(cfg.Database.SQLite.File)
	}
}
