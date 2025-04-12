package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"go.uber.org/zap"
)

// AuthConfig 定义认证中间件的配置
type AuthConfig struct {
	// 排除不需要验证的路径
	ExcludedPaths []string
}

// DefaultAuthConfig 默认认证中间件配置
var DefaultAuthConfig = AuthConfig{
	// 默认排除登录和注册路径
	ExcludedPaths: []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	},
}

// JWTAuth 创建一个JWT认证中间件
func JWTAuth(config ...AuthConfig) fiber.Handler {
	// 使用默认配置或提供的配置
	cfg := DefaultAuthConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		// 检查是否在排除路径中
		path := c.Path()
		for _, excludedPath := range cfg.ExcludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				utils.Debug("跳过JWT认证",
					zap.String("path", path),
					zap.String("excluded_path", excludedPath),
				)
				return c.Next()
			}
		}

		utils.Debug("执行JWT认证", zap.String("path", path))

		// 从请求头获取JWT令牌
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			utils.Warn("未提供认证令牌",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未提供认证令牌")
		}

		// 处理Bearer令牌格式
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 验证令牌
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case utils.ErrTokenExpired:
				utils.Warn("认证令牌已过期",
					zap.String("path", path),
					zap.String("ip", c.IP()),
				)
				return utils.UnauthorizedError(c, "认证令牌已过期")
			case utils.ErrTokenInvalid:
				utils.Warn("认证令牌无效",
					zap.String("path", path),
					zap.String("ip", c.IP()),
				)
				return utils.UnauthorizedError(c, "认证令牌无效")
			case utils.ErrTokenNotProvided:
				utils.Warn("未提供认证令牌",
					zap.String("path", path),
					zap.String("ip", c.IP()),
				)
				return utils.UnauthorizedError(c, "未提供认证令牌")
			default:
				utils.Error("认证失败",
					zap.String("path", path),
					zap.String("ip", c.IP()),
					zap.Error(err),
				)
				return utils.UnauthorizedError(c, "认证失败："+err.Error())
			}
		}

		// 将用户信息存储在上下文中，以便后续处理程序使用
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		utils.Debug("JWT认证成功",
			zap.String("path", path),
			zap.String("ip", c.IP()),
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
		)

		// 继续处理请求
		return c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		utils.Warn("未能从上下文获取用户ID", zap.String("path", c.Path()))
		return 0
	}
	return userID
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *fiber.Ctx) string {
	username, ok := c.Locals("username").(string)
	if !ok {
		utils.Warn("未能从上下文获取用户名", zap.String("path", c.Path()))
		return ""
	}
	return username
}
