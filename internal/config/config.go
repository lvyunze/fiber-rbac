package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port      int    `mapstructure:"port"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"server"`

	Database struct {
		Type   string `mapstructure:"type"`
		SQLite struct {
			File string `mapstructure:"file"`
		} `mapstructure:"sqlite"`
		MySQL struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			DBName   string `mapstructure:"dbname"`
		} `mapstructure:"mysql"`
		Postgres struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			DBName   string `mapstructure:"dbname"`
			SSLMode  string `mapstructure:"sslmode"`
		} `mapstructure:"postgres"`
	} `mapstructure:"database"`

	IPLimit struct {
		Enabled         bool     `mapstructure:"enabled"`
		WhitelistMode   bool     `mapstructure:"whitelist_mode"`
		Whitelist       []string `mapstructure:"whitelist"`
		Blacklist       []string `mapstructure:"blacklist"`
		AllowedNetworks []string `mapstructure:"allowed_networks"`
	} `mapstructure:"ip_limit"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}

// GetDSN 根据配置返回数据库连接字符串
func (c *Config) GetDSN() string {
	switch c.Database.Type {
	case "sqlite":
		return c.Database.SQLite.File
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Database.MySQL.User,
			c.Database.MySQL.Password,
			c.Database.MySQL.Host,
			c.Database.MySQL.Port,
			c.Database.MySQL.DBName,
		)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Postgres.Host,
			c.Database.Postgres.Port,
			c.Database.Postgres.User,
			c.Database.Postgres.Password,
			c.Database.Postgres.DBName,
			c.Database.Postgres.SSLMode,
		)
	default:
		return ""
	}
}
