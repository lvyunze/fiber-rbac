package role

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ListPermissionsHandler 角色权限列表处理器
type ListPermissionsHandler struct {
	roleService service.RoleService
}

// NewListPermissionsHandler 创建角色权限列表处理器
func NewListPermissionsHandler(roleService service.RoleService) *ListPermissionsHandler {
	return &ListPermissionsHandler{
		roleService: roleService,
	}
}

// Handle 处理获取角色权限列表请求
func (h *ListPermissionsHandler) Handle(c *fiber.Ctx) error {
	// 解析角色ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的角色ID")
	}

	// 调用服务层获取角色权限列表
	permissions, err := h.roleService.GetPermissions(id)
	if err != nil {
		slog.Error("获取角色权限列表失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrRoleNotFound {
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		}
		return response.ServerError(c, "获取角色权限列表失败")
	}

	// 返回角色权限列表
	return response.Success(c, permissions, "获取成功")
}
