package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
)

// RegisterSwaggerRoute 封装 swagger UI 路由注册，可通过 enable 参数控制是否开放
func RegisterSwaggerRoute(app *fiber.App, enable bool) {
	if enable {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}
}
