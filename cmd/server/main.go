package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/app"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/logger"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger.Setup(&cfg.Log, cfg.Env)

	// 初始化数据库
	err = model.InitDB(&cfg.Database, cfg.Env)
	if err != nil {
		slog.Error("数据库连接失败", "error", err)
		os.Exit(1)
	}

	// 获取数据库连接
	db := model.GetDB()

	// 自动迁移数据库表结构（服务启动时）
	if err := model.AutoMigrate(db); err != nil {
		slog.Error("数据库自动迁移失败", "error", err)
		os.Exit(1)
	}

	// 初始化默认数据
	if err := model.InitDefaultData(db); err != nil {
		slog.Error("初始化默认数据失败", "error", err)
		os.Exit(1)
	}

	// 初始化验证器
	validator.Init()

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, roleRepo, permissionRepo, refreshTokenRepo, &cfg.JWT)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// 初始化Fiber应用
	fiberApp := app.NewFiberApp(cfg)

	// 注册路由
	app.RegisterRoutes(fiberApp, userService, roleService, permissionService, &cfg.JWT)

	// 启动服务器（非阻塞）
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		slog.Info("服务器启动", "addr", addr)
		if err := fiberApp.Listen(addr); err != nil {
			slog.Error("服务器启动失败", "error", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关闭Fiber应用
	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		slog.Error("服务器关闭失败", "error", err)
	}

	slog.Info("服务器已关闭")
}
