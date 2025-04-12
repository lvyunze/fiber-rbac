package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"go.uber.org/zap"
)

// RequireRole 验证用户是否拥有指定角色
func RequireRole(userService service.UserService, roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		utils.Debug("检查角色权限",
			zap.String("path", path),
			zap.String("role", roleName),
		)

		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			utils.Warn("未能获取用户ID",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定角色
		hasRole, err := userService.HasRole(userID, roleName)
		if err != nil {
			utils.Error("检查角色权限失败",
				zap.String("path", path),
				zap.Uint("user_id", userID),
				zap.String("role", roleName),
				zap.Error(err),
			)
			return utils.ServerError(c, "检查角色失败: "+err.Error())
		}

		if !hasRole {
			utils.Warn("用户没有所需角色",
				zap.String("path", path),
				zap.Uint("user_id", userID),
				zap.String("role", roleName),
			)
			return utils.ForbiddenError(c, "您没有权限访问此资源")
		}

		utils.Debug("角色权限验证通过",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.String("role", roleName),
		)
		return c.Next()
	}
}

// RequirePermission 验证用户是否拥有指定权限
func RequirePermission(userService service.UserService, roleService service.RoleService, permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		utils.Debug("检查权限",
			zap.String("path", path),
			zap.String("permission", permissionName),
		)

		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			utils.Warn("未能获取用户ID",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 获取用户角色
		roles, err := userService.GetUserRoles(userID)
		if err != nil {
			utils.Error("获取用户角色失败",
				zap.String("path", path),
				zap.Uint("user_id", userID),
				zap.Error(err),
			)
			return utils.ServerError(c, "获取用户角色失败: "+err.Error())
		}

		utils.Debug("获取到用户角色",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.Int("role_count", len(roles)),
		)

		// 检查用户的角色是否拥有指定权限
		for _, role := range roles {
			hasPermission, err := roleService.HasPermission(role.ID, permissionName)
			if err != nil {
				utils.Debug("检查角色权限失败",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.Uint("role_id", role.ID),
					zap.String("role_name", role.Name),
					zap.String("permission", permissionName),
					zap.Error(err),
				)
				continue
			}

			if hasPermission {
				utils.Debug("权限验证通过",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.Uint("role_id", role.ID),
					zap.String("role_name", role.Name),
					zap.String("permission", permissionName),
				)
				return c.Next()
			}
		}

		utils.Warn("用户没有所需权限",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.String("permission", permissionName),
		)
		return utils.ForbiddenError(c, "您没有权限执行此操作")
	}
}

// RequireAnyRole 验证用户是否拥有指定角色中的任意一个
func RequireAnyRole(userService service.UserService, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		utils.Debug("检查多角色权限",
			zap.String("path", path),
			zap.Strings("roles", roleNames),
		)

		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			utils.Warn("未能获取用户ID",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定角色中的任意一个
		for _, roleName := range roleNames {
			hasRole, err := userService.HasRole(userID, roleName)
			if err != nil {
				utils.Debug("检查角色权限失败",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.String("role", roleName),
					zap.Error(err),
				)
				continue
			}

			if hasRole {
				utils.Debug("角色权限验证通过",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.String("role", roleName),
				)
				return c.Next()
			}
		}

		utils.Warn("用户没有所需角色",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.Strings("roles", roleNames),
		)
		return utils.ForbiddenError(c, "您没有权限访问此资源")
	}
}

// RequireAllRoles 验证用户是否拥有指定的所有角色
func RequireAllRoles(userService service.UserService, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		utils.Debug("检查全部角色权限",
			zap.String("path", path),
			zap.Strings("roles", roleNames),
		)

		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			utils.Warn("未能获取用户ID",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定的所有角色
		for _, roleName := range roleNames {
			hasRole, err := userService.HasRole(userID, roleName)
			if err != nil {
				utils.Error("检查角色权限失败",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.String("role", roleName),
					zap.Error(err),
				)
				return utils.ServerError(c, "检查角色失败: "+err.Error())
			}

			if !hasRole {
				utils.Warn("用户没有所需角色",
					zap.String("path", path),
					zap.Uint("user_id", userID),
					zap.String("missing_role", roleName),
				)
				return utils.ForbiddenError(c, "您没有权限访问此资源")
			}
		}

		utils.Debug("全部角色权限验证通过",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.Strings("roles", roleNames),
		)
		return c.Next()
	}
}

// RequireAnyPermission 验证用户是否拥有指定权限中的任意一个
func RequireAnyPermission(userService service.UserService, roleService service.RoleService, permissionNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		utils.Debug("检查多权限",
			zap.String("path", path),
			zap.Strings("permissions", permissionNames),
		)

		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			utils.Warn("未能获取用户ID",
				zap.String("path", path),
				zap.String("ip", c.IP()),
			)
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 获取用户角色
		roles, err := userService.GetUserRoles(userID)
		if err != nil {
			utils.Error("获取用户角色失败",
				zap.String("path", path),
				zap.Uint("user_id", userID),
				zap.Error(err),
			)
			return utils.ServerError(c, "获取用户角色失败: "+err.Error())
		}

		utils.Debug("获取到用户角色",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.Int("role_count", len(roles)),
		)

		// 检查用户的角色是否拥有指定权限中的任意一个
		for _, role := range roles {
			for _, permissionName := range permissionNames {
				hasPermission, err := roleService.HasPermission(role.ID, permissionName)
				if err != nil {
					utils.Debug("检查角色权限失败",
						zap.String("path", path),
						zap.Uint("user_id", userID),
						zap.Uint("role_id", role.ID),
						zap.String("role_name", role.Name),
						zap.String("permission", permissionName),
						zap.Error(err),
					)
					continue
				}

				if hasPermission {
					utils.Debug("权限验证通过",
						zap.String("path", path),
						zap.Uint("user_id", userID),
						zap.Uint("role_id", role.ID),
						zap.String("role_name", role.Name),
						zap.String("permission", permissionName),
					)
					return c.Next()
				}
			}
		}

		utils.Warn("用户没有所需权限",
			zap.String("path", path),
			zap.Uint("user_id", userID),
			zap.Strings("permissions", permissionNames),
		)
		return utils.ForbiddenError(c, "您没有权限执行此操作")
	}
}
