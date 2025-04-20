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

// AssignRoleHandler 用户分配角色处理器
type AssignRoleHandler struct {
	userService service.UserService
}

// NewAssignRoleHandler 创建用户分配角色处理器
func NewAssignRoleHandler(userService service.UserService) *AssignRoleHandler {
	return &AssignRoleHandler{
		userService: userService,
	}
}

// Handle 处理用户分配角色请求
// @Summary 分配用户角色
// @Description 为指定用户分配一个或多个角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.AssignRoleRequest true "用户ID与角色ID列表"
// @Success 200 {object} nil "分配成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "用户或角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/assign_role [post]
func (h *AssignRoleHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.AssignRoleRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层分配角色
	err := h.userService.AssignRole(req.UserID, req.RoleIDs)
	if err != nil {
		slog.Error("用户分配角色失败", "userID", req.UserID, "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrUserNotFound:
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		case errors.ErrRoleNotFound:
			return response.Fail(c, response.CodeNotFound, "部分角色不存在")
		default:
			return response.ServerError(c, "用户分配角色失败")
		}
	}

	// 返回分配成功响应
	return response.Success(c, nil, "角色分配成功")
}
