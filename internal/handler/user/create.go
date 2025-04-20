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

// CreateHandler 用户创建处理器
type CreateHandler struct {
	userService service.UserService
}

// NewCreateHandler 创建用户处理器
func NewCreateHandler(userService service.UserService) *CreateHandler {
	return &CreateHandler{
		userService: userService,
	}
}

// Handle 处理创建用户请求
// @Summary 创建用户
// @Description 管理员创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.CreateUserRequest true "用户创建参数"
// @Success 200 {object} map[string]interface{} "创建成功，返回用户ID"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 409 {object} response.Response "用户名或邮箱已存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/create [post]
func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.CreateUserRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层创建用户
	userID, err := h.userService.Create(req)
	if err != nil {
		slog.Error("创建用户失败", "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrUserExists:
			return response.Fail(c, response.CodeParamError, "用户名已存在")
		case errors.ErrEmailExists:
			return response.Fail(c, response.CodeParamError, "邮箱已被使用")
		default:
			return response.ServerError(c, "创建用户失败")
		}
	}

	// 返回创建成功响应
	return response.Success(c, fiber.Map{"id": userID}, "用户创建成功")
}
