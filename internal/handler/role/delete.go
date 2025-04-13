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

// DeleteHandler 角色删除处理器
type DeleteHandler struct {
	roleService service.RoleService
}

// NewDeleteHandler 创建角色删除处理器
func NewDeleteHandler(roleService service.RoleService) *DeleteHandler {
	return &DeleteHandler{
		roleService: roleService,
	}
}

// Handle 处理删除角色请求
func (h *DeleteHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.DeleteRoleRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层删除角色
	err := h.roleService.Delete(req.ID)
	if err != nil {
		slog.Error("删除角色失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrRoleNotFound:
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		case errors.ErrRoleInUse:
			return response.Fail(c, response.CodeForbidden, "角色正在使用中，无法删除")
		default:
			return response.ServerError(c, "删除角色失败")
		}
	}

	// 返回删除成功响应
	return response.Success(c, nil, "角色删除成功")
}
