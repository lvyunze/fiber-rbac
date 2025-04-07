package main

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	// Load configuration
	_ = config.LoadConfig()

	// Initialize database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})

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
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

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
