package user

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

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
	// 解析用户ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的用户ID")
	}

	// 调用服务层删除用户
	err = h.userService.Delete(id)
	if err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "删除用户失败")
	}

	// 返回删除成功响应
	return response.Success(c, nil, "用户删除成功")
}
