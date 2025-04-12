package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// 请求结构
type UserRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

// 注册用户-角色关联管理路由
func RegisterUserRoleRoutes(router fiber.Router, userService service.UserService) {
	router.Post("/:id/roles", assignRolesToUserHandler(userService))
	router.Delete("/:id/roles", removeRolesFromUserHandler(userService))
	router.Get("/:id/roles", getUserRolesHandler(userService))
	router.Get("/:id/has-role/:role", hasRoleHandler(userService))
}

// 分配角色给用户
func assignRolesToUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析用户ID
		userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		// 解析请求体
		var req UserRolesRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestError(c, "无法解析请求体")
		}

		// 检查角色ID列表是否为空
		if len(req.RoleIDs) == 0 {
			return utils.BadRequestError(c, "角色ID列表不能为空")
		}

		// 分配角色
		if err := userService.AssignRolesToUser(uint(userID), req.RoleIDs); err != nil {
			return utils.ServerError(c, "分配角色失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "角色分配成功", nil)
	}
}

// 移除用户的角色
func removeRolesFromUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析用户ID
		userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		// 解析请求体
		var req UserRolesRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestError(c, "无法解析请求体")
		}

		// 检查角色ID列表是否为空
		if len(req.RoleIDs) == 0 {
			return utils.BadRequestError(c, "角色ID列表不能为空")
		}

		// 移除角色
		if err := userService.RemoveRolesFromUser(uint(userID), req.RoleIDs); err != nil {
			return utils.ServerError(c, "移除角色失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "角色移除成功", nil)
	}
}

// 获取用户的所有角色
func getUserRolesHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析用户ID
		userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		// 获取用户角色
		roles, err := userService.GetUserRoles(uint(userID))
		if err != nil {
			return utils.ServerError(c, "获取用户角色失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "获取用户角色成功", roles)
	}
}

// 检查用户是否拥有指定角色
func hasRoleHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析用户ID
		userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		// 获取角色名称
		roleName := c.Params("role")
		if roleName == "" {
			return utils.BadRequestError(c, "角色名称不能为空")
		}

		// 检查用户是否拥有该角色
		hasRole, err := userService.HasRole(uint(userID), roleName)
		if err != nil {
			return utils.ServerError(c, "检查用户角色失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "检查用户角色成功", fiber.Map{
			"has_role": hasRole,
		})
	}
}
