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

// DetailHandler 用户详情处理器
type DetailHandler struct {
	userService service.UserService
}

// NewDetailHandler 创建用户详情处理器
func NewDetailHandler(userService service.UserService) *DetailHandler {
	return &DetailHandler{
		userService: userService,
	}
}

// Handle 处理获取用户详情请求
// @Summary 获取用户详情
// @Description 根据用户ID获取详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.UserDetailRequest true "用户详情参数"
// @Success 200 {object} schema.UserResponse "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/detail [post]
func (h *DetailHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.UserDetailRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// 调用服务层获取用户详情
	user, err := h.userService.GetByID(req.ID)
	if err != nil {
		slog.Error("获取用户详情失败", "id", req.ID, "error", err)
		
		// 处理特定错误类型
		if err == errors.ErrUserNotFound {
			return response.Fail(c, response.CodeNotFound, "用户不存在")
		}
		return response.ServerError(c, "获取用户详情失败")
	}

	// 返回用户详情
	return response.Success(c, user, "获取用户详情成功")
}
