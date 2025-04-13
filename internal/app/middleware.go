package app

import (
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// SetupMiddlewares 设置应用中间件
func SetupMiddlewares(app *fiber.App, cfg *config.Config) {
	// 初始化验证器
	validator.Init()

	// 注册全局中间件
	// 已在fiber.go中注册了基本中间件

	// 注册路由级别中间件
	// 这里可以添加特定路由的中间件
}
