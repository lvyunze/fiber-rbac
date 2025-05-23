package role

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ListHandler 角色列表处理器
type ListHandler struct {
	roleService service.RoleService
}

// NewListHandler 创建角色列表处理器
func NewListHandler(roleService service.RoleService) *ListHandler {
	return &ListHandler{
		roleService: roleService,
	}
}

// Handle 处理获取角色列表请求
// @Summary 获取角色列表
// @Description 分页查询所有角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body schema.ListRoleRequest true "分页与筛选参数"
// @Success 200 {object} schema.ListRoleResponse "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/list [post]
func (h *ListHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.ListRoleRequest)
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

	// 调用服务层获取角色列表
	result, err := h.roleService.List(req)
	if err != nil {
		slog.Error("获取角色列表失败", "error", err)
		return response.ServerError(c, "获取角色列表失败")
	}

	// 返回角色列表
	return response.Success(c, result, "获取成功")
}
