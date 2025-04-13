package auth

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	userService service.UserService
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
	}
}

// Handle 处理登录请求
func (h *LoginHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.LoginRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层进行登录
	res, err := h.userService.Login(req)
	if err != nil {
		slog.Error("用户登录失败", "username", req.Username, "error", err)
		return response.Fail(c, response.CodeUnauthorized, "用户名或密码错误")
	}

	// 返回登录成功响应
	return response.Success(c, res, "登录成功")
}
