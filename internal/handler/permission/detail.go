package permission

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"strconv"

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
	// 解析权限ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Fail(c, response.CodeParamError, "无效的权限ID")
	}

	// 调用服务层获取权限详情
	permission, err := h.permissionService.GetByID(id)
	if err != nil {
		slog.Error("获取权限详情失败", "id", id, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrPermissionNotFound {
			return response.Fail(c, response.CodeNotFound, "权限不存在")
		}
		return response.ServerError(c, "获取权限详情失败")
	}

	// 返回权限详情
	return response.Success(c, permission, "获取成功")
}
