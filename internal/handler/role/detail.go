package role

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// DetailHandler 角色详情处理器
type DetailHandler struct {
	roleService service.RoleService
}

// NewDetailHandler 创建角色详情处理器
func NewDetailHandler(roleService service.RoleService) *DetailHandler {
	return &DetailHandler{
		roleService: roleService,
	}
}

// Handle 处理获取角色详情请求
func (h *DetailHandler) Handle(c *fiber.Ctx) error {
	// 解析角色ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的角色ID")
	}

	// 调用服务层获取角色详情
	role, err := h.roleService.GetByID(id)
	if err != nil {
		slog.Error("获取角色详情失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrRoleNotFound {
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		}
		return response.ServerError(c, "获取角色详情失败")
	}

	// 返回角色详情
	return response.Success(c, role, "获取成功")
}
