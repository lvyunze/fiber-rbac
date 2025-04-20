package user

import (
	"log/slog"

	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ListHandler 用户列表处理器
type ListHandler struct {
	userService service.UserService
}

// NewListHandler 创建用户列表处理器
func NewListHandler(userService service.UserService) *ListHandler {
	return &ListHandler{
		userService: userService,
	}
}

// Handle 处理获取用户列表请求
// @Summary 获取用户列表
// @Description 分页获取用户信息列表，支持关键字搜索
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body schema.ListUserRequest true "分页参数与搜索关键字"
// @Success 200 {object} schema.ListUserResponse
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/list [post]
func (h *ListHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.ListUserRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return response.ParamError(c, err.Error())
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

	// 调用服务层获取用户列表
	result, err := h.userService.List(req)
	if err != nil {
		slog.Error("获取用户列表失败", "error", err)
		return response.ServerError(c, "获取用户列表失败")
	}

	// 返回用户列表
	return response.Success(c, result, "获取成功")
}
