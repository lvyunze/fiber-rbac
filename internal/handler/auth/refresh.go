package auth

import (
	"log/slog"
	"strings"
	
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/lvyunze/fiber-rbac/internal/pkg/validator"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RefreshHandler 刷新令牌处理器
type RefreshHandler struct {
	userService service.UserService
}

// NewRefreshHandler 创建刷新令牌处理器
func NewRefreshHandler(userService service.UserService) *RefreshHandler {
	return &RefreshHandler{
		userService: userService,
	}
}

// Handle 处理刷新令牌请求
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body schema.RefreshTokenRequest true "刷新令牌请求参数"
// @Success 200 {object} schema.LoginResponse "刷新成功，返回新令牌"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "无效的刷新令牌"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *RefreshHandler) Handle(c *fiber.Ctx) error {
	// 解析请求参数
	req := new(schema.RefreshTokenRequest)
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}
	
	var token string
	
	// 优先从请求体获取刷新令牌
	if req.RefreshToken != "" {
		token = req.RefreshToken
	} else {
		// 从 Authorization 头中获取令牌作为备选
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Fail(c, response.CodeUnauthorized, "未提供刷新令牌")
		}
		
		// 解析 Bearer 令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Fail(c, response.CodeUnauthorized, "授权头格式无效")
		}
		
		token = parts[1]
		if token == "" {
			return response.Fail(c, response.CodeUnauthorized, "未提供令牌")
		}
	}

	// 调用服务层刷新令牌
	res, err := h.userService.RefreshToken(token)
	if err != nil {
		slog.Error("刷新令牌失败", "error", err)
		return response.Fail(c, response.CodeUnauthorized, "无效的刷新令牌")
	}

	// 返回刷新成功响应
	return response.Success(c, res, "刷新成功")
}
