package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// RefreshTokenRequest 刷新令牌请求结构
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Token    string       `json:"token"`
	UserInfo *models.User `json:"user_info"`
}

// RegisterAuthRoutes 注册认证相关的路由
func RegisterAuthRoutes(router fiber.Router, userService service.UserService) {
	auth := router.Group("/auth")

	auth.Post("/login", loginHandler(userService))
	auth.Post("/register", registerHandler(userService))
	auth.Post("/refresh", refreshTokenHandler())
}

// loginHandler 处理用户登录
func loginHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析请求体
		req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		// 验证必填字段
		if req.Username == "" || req.Password == "" {
			return utils.BadRequestError(c, "用户名和密码不能为空")
		}

		// 查找用户
		user, err := userService.GetUserByUsername(req.Username)
		if err != nil {
			return utils.UnauthorizedError(c, "用户名或密码错误")
		}

		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return utils.UnauthorizedError(c, "用户名或密码错误")
		}

		// 生成JWT令牌
		token, err := utils.GenerateToken(user.ID, user.Username)
		if err != nil {
			return utils.ServerError(c, "生成令牌失败")
		}

		// 移除敏感信息
		user.Password = ""

		// 返回令牌和用户信息
		return utils.SuccessResponse(c, "登录成功", AuthResponse{
			Token:    token,
			UserInfo: user,
		})
	}
}

// registerHandler 处理用户注册
func registerHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析请求体
		req := new(RegisterRequest)
		if err := c.BodyParser(req); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		// 验证必填字段
		if req.Username == "" || req.Password == "" {
			return utils.BadRequestError(c, "用户名和密码不能为空")
		}

		// 检查密码长度
		if len(req.Password) < 6 {
			return utils.BadRequestError(c, "密码长度必须至少为6个字符")
		}

		// 检查用户名是否已存在
		_, err := userService.GetUserByUsername(req.Username)
		if err == nil {
			return utils.BadRequestError(c, "用户名已存在")
		}

		// 哈希密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return utils.ServerError(c, "密码加密失败")
		}

		// 创建用户
		user := &models.User{
			Username: req.Username,
			Password: string(hashedPassword),
		}

		if err := userService.CreateUser(user); err != nil {
			return utils.ServerError(c, "创建用户失败")
		}

		// 生成JWT令牌
		token, err := utils.GenerateToken(user.ID, user.Username)
		if err != nil {
			return utils.ServerError(c, "生成令牌失败")
		}

		// 移除敏感信息
		user.Password = ""

		// 返回令牌和用户信息
		return utils.SuccessResponse(c, "注册成功", AuthResponse{
			Token:    token,
			UserInfo: user,
		})
	}
}

// refreshTokenHandler 处理令牌刷新
func refreshTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// 首先尝试从请求体中获取令牌
		req := new(RefreshTokenRequest)
		if err := c.BodyParser(req); err == nil && req.Token != "" {
			tokenString = req.Token
		} else {
			// 从请求头获取JWT令牌
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return utils.TokenError(c, "未提供令牌", utils.CodeTokenNotProvided)
			}

			// 处理Bearer令牌格式
			tokenString = authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 验证令牌格式
		if tokenString == "" {
			return utils.TokenError(c, "未提供令牌", utils.CodeTokenNotProvided)
		}

		// 刷新令牌
		newToken, err := utils.RefreshToken(tokenString)
		if err != nil {
			if err == utils.ErrTokenInvalid {
				return utils.TokenError(c, "无效的令牌", utils.CodeTokenInvalid)
			} else if err == utils.ErrTokenExpired {
				// 即使令牌过期，也尝试刷新
				claims, _ := utils.ExtractClaims(tokenString)
				if claims != nil {
					newToken, err = utils.GenerateToken(claims.UserID, claims.Username)
					if err == nil {
						return utils.SuccessResponse(c, "刷新令牌成功", fiber.Map{
							"token": newToken,
						})
					}
				}
				return utils.TokenError(c, "令牌已过期", utils.CodeTokenExpired)
			}
			return utils.ServerError(c, "刷新令牌失败")
		}

		// 返回新令牌
		return utils.SuccessResponse(c, "刷新令牌成功", fiber.Map{
			"token": newToken,
		})
	}
}
