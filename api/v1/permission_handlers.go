package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// RegisterPermissionRoutes 注册权限相关的路由
func RegisterPermissionRoutes(router fiber.Router, permissionService service.PermissionService) {
	permissions := router.Group("/permissions")

	permissions.Post("/", createPermissionHandler(permissionService))
	permissions.Get("/", getPermissionsHandler(permissionService))
	permissions.Get("/:id", getPermissionByIDHandler(permissionService))
	permissions.Put("/:id", updatePermissionByIDHandler(permissionService))
	permissions.Delete("/:id", deletePermissionByIDHandler(permissionService))
}

func createPermissionHandler(permissionService service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permission := new(models.Permission)
		if err := c.BodyParser(permission); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := permissionService.CreatePermission(permission); err != nil {
			return utils.ServerError(c, "创建权限失败")
		}

		return utils.SuccessResponse(c, "权限创建成功", permission)
	}
}

func getPermissionsHandler(permissionService service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, err := permissionService.GetPermissions()
		if err != nil {
			return utils.ServerError(c, "获取权限列表失败")
		}

		return utils.SuccessResponse(c, "获取权限列表成功", permissions)
	}
}

func getPermissionByIDHandler(permissionService service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的权限ID")
		}

		permission, err := permissionService.GetPermissionByID(uint(id))
		if err != nil {
			return utils.NotFoundError(c, "权限不存在", utils.ErrPermissionNotFound)
		}

		return utils.SuccessResponse(c, "获取权限成功", permission)
	}
}

func updatePermissionByIDHandler(permissionService service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的权限ID")
		}

		permission := new(models.Permission)
		if err := c.BodyParser(permission); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := permissionService.UpdatePermissionByID(uint(id), permission); err != nil {
			return utils.ServerError(c, "更新权限失败")
		}

		return utils.SuccessResponse(c, "更新权限成功", permission)
	}
}

func deletePermissionByIDHandler(permissionService service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的权限ID")
		}

		if err := permissionService.DeletePermissionByID(uint(id)); err != nil {
			return utils.ServerError(c, "删除权限失败")
		}

		return utils.SuccessResponse(c, "删除权限成功", nil)
	}
}
