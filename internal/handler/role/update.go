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

// UpdateHandler 角色更新处理器
type UpdateHandler struct {
	roleService service.RoleService
}

// NewUpdateHandler 创建角色更新处理器
func NewUpdateHandler(roleService service.RoleService) *UpdateHandler {
	return &UpdateHandler{
		roleService: roleService,
	}
}

// Handle 处理更新角色请求
func (h *UpdateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UpdateRoleRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层更新角色
	err := h.roleService.Update(req)
	if err != nil {
		slog.Error("更新角色失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrRoleNotFound:
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		case errors.ErrRoleExists:
			return response.Fail(c, response.CodeParamError, "角色名已存在")
		default:
			return response.ServerError(c, "更新角色失败")
		}
	}

	// 返回更新成功响应
	return response.Success(c, nil, "角色更新成功")
}
