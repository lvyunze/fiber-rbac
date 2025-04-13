package user

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ListRolesHandler 用户角色列表处理器
type ListRolesHandler struct {
	userService service.UserService
}

// NewListRolesHandler 创建用户角色列表处理器
func NewListRolesHandler(userService service.UserService) *ListRolesHandler {
	return &ListRolesHandler{
		userService: userService,
	}
}

// Handle 处理获取用户角色列表请求
func (h *ListRolesHandler) Handle(c *fiber.Ctx) error {
	// 解析用户ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的用户ID")
	}

	// 调用服务层获取用户角色列表
	roles, err := h.userService.GetRoles(id)
	if err != nil {
		slog.Error("获取用户角色列表失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "获取用户角色列表失败")
	}

	// 返回用户角色列表
	return response.Success(c, roles, "获取成功")
}
