package auth

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RefreshHandler 刷新令牌处理器
type RefreshHandler struct {
	userService service.UserService
}

// NewRefreshHandler 创建刷新令牌处理器
func NewRefreshHandler(userService service.UserService) *RefreshHandler {
	return &RefreshHandler{
		userService: userService,
	}
}

// Handle 处理刷新令牌请求
func (h *RefreshHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.RefreshTokenRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层刷新令牌
	res, err := h.userService.RefreshToken(req)
	if err != nil {
		slog.Error("刷新令牌失败", "error", err)
		return response.Fail(c, response.CodeUnauthorized, "无效的刷新令牌")
	}

	// 返回刷新成功响应
	return response.Success(c, res, "刷新成功")
}
