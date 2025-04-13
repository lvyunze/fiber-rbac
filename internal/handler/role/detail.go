package role

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

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
	// 解析请求参数
	req := new(schema.RoleDetailRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层获取角色详情
	role, err := h.roleService.GetByID(req.ID)
	if err != nil {
		slog.Error("获取角色详情失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrRoleNotFound {
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		}
		return response.ServerError(c, "获取角色详情失败")
	}

	// 返回角色详情
	return response.Success(c, role, "获取成功")
}
