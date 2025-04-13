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

// CreateHandler 权限创建处理器
type CreateHandler struct {
	permissionService service.PermissionService
}

// NewCreateHandler 创建权限处理器
func NewCreateHandler(permissionService service.PermissionService) *CreateHandler {
	return &CreateHandler{
		permissionService: permissionService,
	}
}

// Handle 处理创建权限请求
func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.CreatePermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层创建权限
	permissionID, err := h.permissionService.Create(req)
	if err != nil {
		slog.Error("创建权限失败", "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrPermissionExists {
			return response.Fail(c, response.CodeParamError, "权限标识已存在")
		}
		return response.ServerError(c, "创建权限失败")
	}

	// 返回创建成功响应
	return response.Success(c, fiber.Map{"id": permissionID}, "权限创建成功")
}
