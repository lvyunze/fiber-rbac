package app

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
)

// NewFiberApp 创建并配置Fiber应用实例
func NewFiberApp(cfg *config.Config) *fiber.App {
	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		AppName:               "RBAC API",
		ReadTimeout:           time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout:          time.Duration(cfg.Server.Timeout) * time.Second,
		IdleTimeout:           time.Duration(cfg.Server.Timeout) * time.Second,
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
		ErrorHandler:          customErrorHandler,
	})

	// 注册中间件
	registerMiddlewares(app, cfg)

	slog.Info("Fiber应用初始化完成")
	return app
}

// registerMiddlewares 注册全局中间件
func registerMiddlewares(app *fiber.App, cfg *config.Config) {
	// 恢复中间件，用于捕获panic
	app.Use(recover.New())

	// IP白名单中间件
	app.Use(middleware.IPWhitelist(&cfg.Security))

	// CORS中间件，允许跨域请求
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Range",
	}))

	// 日志中间件
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Shanghai",
	}))
}

// customErrorHandler 自定义错误处理器
func customErrorHandler(c *fiber.Ctx, err error) error {
	// 默认状态码为500
	code := fiber.StatusInternalServerError

	// 检查是否为Fiber的错误类型
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// 记录错误日志
	slog.Error("请求处理错误", "path", c.Path(), "error", err.Error())

	// 返回JSON格式的错误响应
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"code":    code,
		"message": err.Error(),
		"data":    nil,
	})
}
