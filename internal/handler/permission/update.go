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

// UpdateHandler 权限更新处理器
type UpdateHandler struct {
	permissionService service.PermissionService
}

// NewUpdateHandler 创建权限更新处理器
func NewUpdateHandler(permissionService service.PermissionService) *UpdateHandler {
	return &UpdateHandler{
		permissionService: permissionService,
	}
}

// Handle 处理更新权限请求
func (h *UpdateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UpdatePermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层更新权限
	err := h.permissionService.Update(req)
	if err != nil {
		slog.Error("更新权限失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrPermissionNotFound:
			return response.Fail(c, response.CodeNotFound, "权限不存在")
		case errors.ErrPermissionExists:
			return response.Fail(c, response.CodeParamError, "权限标识已存在")
		default:
			return response.ServerError(c, "更新权限失败")
		}
	}

	// 返回更新成功响应
	return response.Success(c, nil, "权限更新成功")
}
