package auth

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// CheckHandler 权限检查处理器
type CheckHandler struct {
	userService service.UserService
}

// NewCheckHandler 创建权限检查处理器
func NewCheckHandler(userService service.UserService) *CheckHandler {
	return &CheckHandler{
		userService: userService,
	}
}

// Handle 处理权限检查请求
func (h *CheckHandler) Handle(c *fiber.Ctx) error {
	// 从上下文获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "无效的授权令牌")
	}

	// 解析请求参数
	req := new(schema.CheckPermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 检查用户权限
	hasPermission, err := h.userService.CheckPermission(userID, req.Permission)
	if err != nil {
		slog.Error("检查权限失败", "userID", userID, "permission", req.Permission, "error", err)
		return response.ServerError(c, "检查权限失败")
	}

	// 返回检查结果
	if hasPermission {
		return response.Success(c, fiber.Map{"has_permission": true}, "用户具有该权限")
	} else {
		return response.Success(c, fiber.Map{"has_permission": false}, "用户不具有该权限")
	}
}
