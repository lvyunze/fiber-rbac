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

// CreateHandler 角色创建处理器
type CreateHandler struct {
	roleService service.RoleService
}

// NewCreateHandler 创建角色处理器
func NewCreateHandler(roleService service.RoleService) *CreateHandler {
	return &CreateHandler{
		roleService: roleService,
	}
}

// Handle 处理创建角色请求
func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.CreateRoleRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层创建角色
	roleID, err := h.roleService.Create(req)
	if err != nil {
		slog.Error("创建角色失败", "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrRoleExists {
			return response.Fail(c, response.CodeParamError, "角色名已存在")
		}
		return response.ServerError(c, "创建角色失败")
	}

	// 返回创建成功响应
	return response.Success(c, fiber.Map{"id": roleID}, "角色创建成功")
}
