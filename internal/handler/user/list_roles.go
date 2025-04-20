package user

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ListRolesHandler 用户角色列表处理器
type ListRolesHandler struct {
	userService service.UserService
}

// NewListRolesHandler 创建用户角色列表处理器
func NewListRolesHandler(userService service.UserService) *ListRolesHandler {
	return &ListRolesHandler{
		userService: userService,
	}
}

// Handle 处理获取用户角色列表请求
// @Summary 获取用户角色列表
// @Description 查询指定用户已分配的所有角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.UserListRolesRequest true "用户ID参数"
// @Success 200 {object} []schema.RoleSimple "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/list_roles [post]
func (h *ListRolesHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UserListRolesRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层获取用户角色列表
	roles, err := h.userService.GetRoles(req.ID)
	if err != nil {
		slog.Error("获取用户角色列表失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "获取用户角色列表失败")
	}

	// 返回用户角色列表
	return response.Success(c, roles, "获取成功")
}
