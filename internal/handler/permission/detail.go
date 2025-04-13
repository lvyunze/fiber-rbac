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

// DetailHandler 权限详情处理器
type DetailHandler struct {
	permissionService service.PermissionService
}

// NewDetailHandler 创建权限详情处理器
func NewDetailHandler(permissionService service.PermissionService) *DetailHandler {
	return &DetailHandler{
		permissionService: permissionService,
	}
}

// Handle 处理获取权限详情请求
func (h *DetailHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.GetPermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层获取权限详情
	permission, err := h.permissionService.GetByID(req.ID)
	if err != nil {
		slog.Error("获取权限详情失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrPermissionNotFound {
			return response.Fail(c, response.CodeNotFound, "权限不存在")
		}
		return response.ServerError(c, "获取权限详情失败")
	}

	// 返回权限详情
	return response.Success(c, permission, "获取成功")
}
