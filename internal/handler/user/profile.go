package user

import (
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ProfileHandler 用户个人信息处理器
type ProfileHandler struct {
	userService service.UserService
}

// NewProfileHandler 创建个人信息处理器
func NewProfileHandler(userService service.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// Handle 处理获取个人信息请求
func (h *ProfileHandler) Handle(c *fiber.Ctx) error {
	// 从上下文中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "未授权的访问")
	}

	// 调用服务层获取用户信息
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return response.Fail(c, response.CodeServerError, "获取用户信息失败")
	}

	// 返回用户信息
	return response.Success(c, user, "获取成功")
}
