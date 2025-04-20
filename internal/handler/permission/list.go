package permission

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ListHandler 权限列表处理器
type ListHandler struct {
	permissionService service.PermissionService
}

// NewListHandler 创建权限列表处理器
func NewListHandler(permissionService service.PermissionService) *ListHandler {
	return &ListHandler{
		permissionService: permissionService,
	}
}

// Handle 处理获取权限列表请求
// @Summary 获取权限列表
// @Description 分页查询所有权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param data body schema.ListPermissionRequest true "分页与筛选参数"
// @Success 200 {object} schema.ListPermissionResponse "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/list [post]
func (h *ListHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.ListPermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	} else if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大每页数量
	}

	// 调用服务层获取权限列表
	result, err := h.permissionService.List(req)
	if err != nil {
		slog.Error("获取权限列表失败", "error", err)
		return response.ServerError(c, "获取权限列表失败")
	}

	// 返回权限列表
	return response.Success(c, result, "获取成功")
}
