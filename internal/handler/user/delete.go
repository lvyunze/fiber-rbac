package user

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// DeleteHandler 用户删除处理器
type DeleteHandler struct {
	userService service.UserService
}

// NewDeleteHandler 创建用户删除处理器
func NewDeleteHandler(userService service.UserService) *DeleteHandler {
	return &DeleteHandler{
		userService: userService,
	}
}

// Handle 处理删除用户请求
func (h *DeleteHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UserDeleteRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 获取当前登录用户ID
	currentUserID, ok := c.Locals("userID").(uint64)
	if !ok || currentUserID == 0 {
		return response.Fail(c, response.CodeUnauthorized, "未登录或身份异常")
	}
	if req.ID == currentUserID {
		return response.Fail(c, response.CodeForbidden, "禁止删除自己的账号")
	}

	// 调用服务层删除用户
	err := h.userService.Delete(req.ID)
	if err != nil {
		slog.Error("删除用户失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "删除用户失败")
	}

	// 返回删除成功响应
	return response.Success(c, nil, "用户删除成功")
}
