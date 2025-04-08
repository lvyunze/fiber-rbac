package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/database"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		panic("无法获取数据库连接: " + err.Error())
	}

	// 测试数据库连接是否正常
	if err := sqlDB.Ping(); err != nil {
		panic("数据库连接测试失败: " + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 输出连接成功信息
	fmt.Println("数据库连接测试成功，类型:", cfg.Database.Type)

	// 配置IP限制中间件
	if viper.GetBool("ip_limit.enabled") {
		whitelistMode := viper.GetBool("ip_limit.whitelist_mode")
		whitelist := viper.GetStringSlice("ip_limit.whitelist")
		blacklist := viper.GetStringSlice("ip_limit.blacklist")
		allowedNetworks := viper.GetStringSlice("ip_limit.allowed_networks")

		ipLimiter := middleware.NewIPLimiter(whitelistMode, whitelist, blacklist, allowedNetworks)
		app.Use(ipLimiter.Handler())
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(database.DB)
	roleRepo := repository.NewRoleRepository(database.DB)
	permissionRepo := repository.NewPermissionRepository(database.DB)

	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// 创建API路由组
	router := app.Group("/api/v1")

	// 注册认证路由（不需要JWT保护）
	v1.RegisterAuthRoutes(router, userService)

	// 配置JWT中间件
	jwtMiddleware := middleware.JWTAuth(middleware.AuthConfig{
		ExcludedPaths: []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/refresh",
		},
	})

	// 应用JWT中间件到所有受保护的路由
	router.Use(jwtMiddleware)

	// 注册其他路由（JWT保护）
	v1.RegisterUserRoutes(router, userService)
	v1.RegisterRoleRoutes(router, roleService)
	v1.RegisterPermissionRoutes(router, permissionService)

	// Start server
	app.Listen(":3000")
}
