package user

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// DetailHandler 用户详情处理器
type DetailHandler struct {
	userService service.UserService
}

// NewDetailHandler 创建用户详情处理器
func NewDetailHandler(userService service.UserService) *DetailHandler {
	return &DetailHandler{
		userService: userService,
	}
}

// Handle 处理获取用户详情请求
func (h *DetailHandler) Handle(c *fiber.Ctx) error {
	// 解析用户ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的用户ID")
	}

	// 调用服务层获取用户详情
	user, err := h.userService.GetByID(id)
	if err != nil {
		slog.Error("获取用户详情失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "获取用户详情失败")
	}

	// 返回用户详情
	return response.Success(c, user, "获取用户详情成功")
}
