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

// UpdateHandler 用户更新处理器
type UpdateHandler struct {
	userService service.UserService
}

// NewUpdateHandler 创建用户更新处理器
func NewUpdateHandler(userService service.UserService) *UpdateHandler {
	return &UpdateHandler{
		userService: userService,
	}
}

// Handle 处理更新用户请求
// @Summary 更新用户信息
// @Description 更新指定用户的基本信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.UpdateUserRequest true "用户更新参数"
// @Success 200 {object} nil "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 409 {object} response.Response "用户名或邮箱已存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/update [post]
func (h *UpdateHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UpdateUserRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层更新用户
	err := h.userService.Update(req)
	if err != nil {
		slog.Error("更新用户失败", "error", err)
		
		// 处理特定错误类型
		switch err {
		case errors.ErrUserNotFound:
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		case errors.ErrUserExists:
			return response.Fail(c, response.CodeParamError, "用户名已存在")
		case errors.ErrEmailExists:
			return response.Fail(c, response.CodeParamError, "邮箱已被使用")
		default:
			return response.ServerError(c, "更新用户失败")
		}
	}

	// 返回更新成功响应
	return response.Success(c, nil, "用户更新成功")
}
