package middleware

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/pkg/jwt"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Auth 认证中间件
func Auth(jwtConfig *config.JWTConfig) fiber.Handler {
	tokenService := jwt.NewTokenService(jwtConfig)

	return func(c *fiber.Ctx) error {
		// 从请求头获取Token
		token := c.Get("Authorization")
		if token == "" {
			return response.Unauthorized(c, "未提供认证令牌")
		}

		// 移除Bearer前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		} else {
			return response.Unauthorized(c, "无效的认证令牌格式")
		}

		// 验证Token
		claims, err := tokenService.ValidateToken(token)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				return response.Unauthorized(c, "认证令牌已过期")
			}
			slog.Error("验证令牌失败", "error", err)
			return response.Unauthorized(c, "无效的认证令牌")
		}

		// 检查令牌类型
		if claims.TokenType != "access" {
			return response.Unauthorized(c, "令牌类型错误")
		}

		// 将用户信息存储到上下文中
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *fiber.Ctx) uint64 {
	userID, ok := c.Locals("userID").(uint64)
	if !ok {
		return 0
	}
	return userID
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *fiber.Ctx) string {
	username, ok := c.Locals("username").(string)
	if !ok {
		return ""
	}
	return username
}
