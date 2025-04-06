package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
)

// RegisterPermissionRoutes 注册权限相关的路由
func RegisterPermissionRoutes(router fiber.Router, permissionService *service.PermissionService) {
	permissions := router.Group("/permissions")

	permissions.Post("/", createPermissionHandler(permissionService))
	permissions.Get("/", getPermissionsHandler(permissionService))
	permissions.Get("/:id", getPermissionByIDHandler(permissionService))
	permissions.Put("/:id", updatePermissionByIDHandler(permissionService))
	permissions.Delete("/:id", deletePermissionByIDHandler(permissionService))
}

func createPermissionHandler(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permission := new(models.Permission)
		if err := c.BodyParser(permission); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := permissionService.CreatePermission(permission); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "创建权限失败",
				"data":    nil,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "权限创建成功",
			"data":    permission,
		})
	}
}

func getPermissionsHandler(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, err := permissionService.GetPermissions()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "获取权限列表失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取权限列表成功",
			"data":    permissions,
		})
	}
}

func getPermissionByIDHandler(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的权限ID",
				"data":    nil,
			})
		}

		permission, err := permissionService.GetPermissionByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "权限不存在",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取权限成功",
			"data":    permission,
		})
	}
}

func updatePermissionByIDHandler(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的权限ID",
				"data":    nil,
			})
		}

		permission := new(models.Permission)
		if err := c.BodyParser(permission); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := permissionService.UpdatePermissionByID(uint(id), permission); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "更新权限失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "更新权限成功",
			"data":    permission,
		})
	}
}

func deletePermissionByIDHandler(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的权限ID",
				"data":    nil,
			})
		}

		if err := permissionService.DeletePermissionByID(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "删除权限失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "删除权限成功",
			"data":    nil,
		})
	}
}
