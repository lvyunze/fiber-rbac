package permission

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// DeleteHandler 权限删除处理器
type DeleteHandler struct {
	permissionService service.PermissionService
}

// NewDeleteHandler 创建权限删除处理器
func NewDeleteHandler(permissionService service.PermissionService) *DeleteHandler {
	return &DeleteHandler{
		permissionService: permissionService,
	}
}

// Handle 处理删除权限请求
func (h *DeleteHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.DeletePermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层删除权限
	err := h.permissionService.Delete(req.ID)
	if err != nil {
		slog.Error("删除权限失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrPermissionNotFound:
			return response.Fail(c, response.CodeNotFound, "权限不存在")
		case errors.ErrPermissionInUse:
			return response.Fail(c, response.CodeForbidden, "权限正在使用中，无法删除")
		default:
			return response.ServerError(c, "删除权限失败")
		}
	}

	// 返回删除成功响应
	return response.Success(c, nil, "权限删除成功")
}
