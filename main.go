package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/database"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		// 错误处理
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// 记录错误
			utils.Error("服务器错误", zap.Error(err), zap.String("path", c.Path()))

			// 默认返回500错误
			code := fiber.StatusInternalServerError

			// 如果错误是*fiber.Error，则使用其状态码
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志系统
	isDevelopment := os.Getenv("APP_ENV") != "production"
	utils.InitLogger(cfg.Logger.Level, cfg.Logger.FilePath, isDevelopment)
	defer utils.Close() // 确保程序退出前日志被正确刷新

	utils.Info("应用启动", zap.String("环境", map[bool]string{true: "开发环境", false: "生产环境"}[isDevelopment]))

	// 添加全局日志中间件
	app.Use(middleware.Logger(middleware.LoggerConfig{
		LogRequestBody:  isDevelopment,                   // 开发环境记录请求体
		LogResponseBody: false,                           // 默认不记录响应体
		SkipPaths:       []string{"/health", "/metrics"}, // 跳过健康检查和指标接口
	}))

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		utils.Fatal("数据库连接失败", zap.Error(err))
	}

	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		utils.Fatal("无法获取数据库连接", zap.Error(err))
	}

	// 测试数据库连接是否正常
	if err := sqlDB.Ping(); err != nil {
		utils.Fatal("数据库连接测试失败", zap.Error(err))
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 输出连接成功信息
	utils.Info("数据库连接成功",
		zap.String("类型", cfg.Database.Type),
		zap.Int("最大空闲连接", 10),
		zap.Int("最大打开连接", 100),
	)

	// 配置IP限制中间件
	if viper.GetBool("ip_limit.enabled") {
		whitelistMode := viper.GetBool("ip_limit.whitelist_mode")
		whitelist := viper.GetStringSlice("ip_limit.whitelist")
		blacklist := viper.GetStringSlice("ip_limit.blacklist")
		allowedNetworks := viper.GetStringSlice("ip_limit.allowed_networks")

		utils.Info("启用IP限制中间件",
			zap.Bool("白名单模式", whitelistMode),
			zap.Strings("白名单", whitelist),
			zap.Strings("黑名单", blacklist),
		)

		ipLimiter := middleware.NewIPLimiter(whitelistMode, whitelist, blacklist, allowedNetworks)
		app.Use(ipLimiter.Handler())
	}

	// 添加健康检查路由
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

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
	utils.Info("注册认证路由", zap.String("路径前缀", "/api/v1/auth"))

	// 配置JWT中间件
	excludedPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
	}

	jwtMiddleware := middleware.JWTAuth(middleware.AuthConfig{
		ExcludedPaths: excludedPaths,
	})
	utils.Info("配置JWT中间件", zap.Strings("排除路径", excludedPaths))

	// 应用JWT中间件到所有受保护的路由
	router.Use(jwtMiddleware)
	utils.Info("应用JWT中间件到API路由")

	// 注册其他路由（JWT保护）
	v1.RegisterUserRoutes(router, userService)
	v1.RegisterRoleRoutes(router, roleService)
	v1.RegisterPermissionRoutes(router, permissionService)
	utils.Info("注册受保护路由",
		zap.String("用户路由", "/api/v1/users"),
		zap.String("角色路由", "/api/v1/roles"),
		zap.String("权限路由", "/api/v1/permissions"),
	)

	// 获取服务器端口
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.Info("启动服务器", zap.String("监听地址", port))

	// 启动服务器
	if err := app.Listen(port); err != nil {
		utils.Fatal("服务器启动失败", zap.Error(err))
	}
}
