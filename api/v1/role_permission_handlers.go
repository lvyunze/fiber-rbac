package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// 请求结构
type RolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

// 注册角色-权限关联管理路由
func RegisterRolePermissionRoutes(router fiber.Router, roleService service.RoleService) {
	router.Post("/:id/permissions", assignPermissionsToRoleHandler(roleService))
	router.Delete("/:id/permissions", removePermissionsFromRoleHandler(roleService))
	router.Get("/:id/permissions", getRolePermissionsHandler(roleService))
	router.Get("/:id/has-permission/:permission", hasPermissionHandler(roleService))
}

// 分配权限给角色
func assignPermissionsToRoleHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析角色ID
		roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		// 解析请求体
		var req RolePermissionsRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestError(c, "无法解析请求体")
		}

		// 检查权限ID列表是否为空
		if len(req.PermissionIDs) == 0 {
			return utils.BadRequestError(c, "权限ID列表不能为空")
		}

		// 分配权限
		if err := roleService.AssignPermissionsToRole(uint(roleID), req.PermissionIDs); err != nil {
			return utils.ServerError(c, "分配权限失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "权限分配成功", nil)
	}
}

// 移除角色的权限
func removePermissionsFromRoleHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析角色ID
		roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		// 解析请求体
		var req RolePermissionsRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestError(c, "无法解析请求体")
		}

		// 检查权限ID列表是否为空
		if len(req.PermissionIDs) == 0 {
			return utils.BadRequestError(c, "权限ID列表不能为空")
		}

		// 移除权限
		if err := roleService.RemovePermissionsFromRole(uint(roleID), req.PermissionIDs); err != nil {
			return utils.ServerError(c, "移除权限失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "权限移除成功", nil)
	}
}

// 获取角色的所有权限
func getRolePermissionsHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析角色ID
		roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		// 获取角色权限
		permissions, err := roleService.GetRolePermissions(uint(roleID))
		if err != nil {
			return utils.ServerError(c, "获取角色权限失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "获取角色权限成功", permissions)
	}
}

// 检查角色是否拥有指定权限
func hasPermissionHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析角色ID
		roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		// 获取权限名称
		permissionName := c.Params("permission")
		if permissionName == "" {
			return utils.BadRequestError(c, "权限名称不能为空")
		}

		// 检查角色是否拥有该权限
		hasPermission, err := roleService.HasPermission(uint(roleID), permissionName)
		if err != nil {
			return utils.ServerError(c, "检查角色权限失败: "+err.Error())
		}

		return utils.SuccessResponse(c, "检查角色权限成功", fiber.Map{
			"has_permission": hasPermission,
		})
	}
}
