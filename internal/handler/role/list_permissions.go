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

// ListPermissionsHandler 角色权限列表处理器
type ListPermissionsHandler struct {
	roleService service.RoleService
}

// NewListPermissionsHandler 创建角色权限列表处理器
func NewListPermissionsHandler(roleService service.RoleService) *ListPermissionsHandler {
	return &ListPermissionsHandler{
		roleService: roleService,
	}
}

// Handle 处理获取角色权限列表请求
// @Summary 获取角色权限列表
// @Description 查询指定角色已分配的所有权限
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param data body schema.RoleListPermissionsRequest true "角色ID参数"
// @Success 200 {object} []schema.PermissionSimple "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/list_permissions [post]
func (h *ListPermissionsHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.RoleListPermissionsRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层获取角色权限列表
	permissions, err := h.roleService.GetPermissions(req.ID)
	if err != nil {
		slog.Error("获取角色权限列表失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrRoleNotFound {
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		}
		return response.ServerError(c, "获取角色权限列表失败")
	}

	// 返回角色权限列表
	return response.Success(c, permissions, "获取成功")
}
