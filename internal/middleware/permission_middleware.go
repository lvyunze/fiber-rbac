package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// RequireRole 验证用户是否拥有指定角色
func RequireRole(userService service.UserService, roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定角色
		hasRole, err := userService.HasRole(userID, roleName)
		if err != nil {
			return utils.ServerError(c, "检查角色失败: "+err.Error())
		}

		if !hasRole {
			return utils.ForbiddenError(c, "您没有权限访问此资源")
		}

		return c.Next()
	}
}

// RequirePermission 验证用户是否拥有指定权限
func RequirePermission(userService service.UserService, roleService service.RoleService, permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 获取用户角色
		roles, err := userService.GetUserRoles(userID)
		if err != nil {
			return utils.ServerError(c, "获取用户角色失败: "+err.Error())
		}

		// 检查用户的角色是否拥有指定权限
		for _, role := range roles {
			hasPermission, err := roleService.HasPermission(role.ID, permissionName)
			if err != nil {
				continue
			}

			if hasPermission {
				return c.Next()
			}
		}

		return utils.ForbiddenError(c, "您没有权限执行此操作")
	}
}

// RequireAnyRole 验证用户是否拥有指定角色中的任意一个
func RequireAnyRole(userService service.UserService, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定角色中的任意一个
		for _, roleName := range roleNames {
			hasRole, err := userService.HasRole(userID, roleName)
			if err != nil {
				continue
			}

			if hasRole {
				return c.Next()
			}
		}

		return utils.ForbiddenError(c, "您没有权限访问此资源")
	}
}

// RequireAllRoles 验证用户是否拥有指定的所有角色
func RequireAllRoles(userService service.UserService, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 检查用户是否拥有指定的所有角色
		for _, roleName := range roleNames {
			hasRole, err := userService.HasRole(userID, roleName)
			if err != nil {
				return utils.ServerError(c, "检查角色失败: "+err.Error())
			}

			if !hasRole {
				return utils.ForbiddenError(c, "您没有权限访问此资源")
			}
		}

		return c.Next()
	}
}

// RequireAnyPermission 验证用户是否拥有指定权限中的任意一个
func RequireAnyPermission(userService service.UserService, roleService service.RoleService, permissionNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文中获取用户ID
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return utils.UnauthorizedError(c, "未登录或会话已过期")
		}

		// 获取用户角色
		roles, err := userService.GetUserRoles(userID)
		if err != nil {
			return utils.ServerError(c, "获取用户角色失败: "+err.Error())
		}

		// 检查用户的角色是否拥有指定权限中的任意一个
		for _, role := range roles {
			for _, permissionName := range permissionNames {
				hasPermission, err := roleService.HasPermission(role.ID, permissionName)
				if err != nil {
					continue
				}

				if hasPermission {
					return c.Next()
				}
			}
		}

		return utils.ForbiddenError(c, "您没有权限执行此操作")
	}
}
