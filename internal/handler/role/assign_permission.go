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

// AssignPermissionHandler 角色分配权限处理器
type AssignPermissionHandler struct {
	roleService service.RoleService
}

// NewAssignPermissionHandler 创建角色分配权限处理器
func NewAssignPermissionHandler(roleService service.RoleService) *AssignPermissionHandler {
	return &AssignPermissionHandler{
		roleService: roleService,
	}
}

// Handle 处理角色分配权限请求
// @Summary 角色分配权限
// @Description 为指定角色分配一个或多个权限
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param data body schema.AssignPermissionRequest true "角色ID与权限ID列表"
// @Success 200 {object} nil "分配成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "角色或权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/assign_permission [post]
func (h *AssignPermissionHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.AssignPermissionRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层分配权限
	err := h.roleService.AssignPermission(req.RoleID, req.PermissionIDs)
	if err != nil {
		slog.Error("角色分配权限失败", "roleID", req.RoleID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrRoleNotFound:
			return response.Fail(c, response.CodeNotFound, "角色不存在")
		case errors.ErrPermissionNotFound:
			return response.Fail(c, response.CodeNotFound, "部分权限不存在")
		default:
			return response.ServerError(c, "角色分配权限失败")
		}
	}

	// 返回分配成功响应
	return response.Success(c, nil, "权限分配成功")
}
