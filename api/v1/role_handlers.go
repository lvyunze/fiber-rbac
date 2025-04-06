package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
)

// RegisterRoleRoutes 注册角色相关的路由
func RegisterRoleRoutes(router fiber.Router, roleService *service.RoleService) {
	roles := router.Group("/roles")

	roles.Post("/", createRoleHandler(roleService))
	roles.Get("/", getRolesHandler(roleService))
	roles.Get("/:id", getRoleByIDHandler(roleService))
	roles.Put("/:id", updateRoleByIDHandler(roleService))
	roles.Delete("/:id", deleteRoleByIDHandler(roleService))
}

func createRoleHandler(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := new(models.Role)
		if err := c.BodyParser(role); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := roleService.CreateRole(role); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "创建角色失败",
				"data":    nil,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "角色创建成功",
			"data":    role,
		})
	}
}

func getRolesHandler(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := roleService.GetRoles()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "获取角色列表失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取角色列表成功",
			"data":    roles,
		})
	}
}

func getRoleByIDHandler(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的角色ID",
				"data":    nil,
			})
		}

		role, err := roleService.GetRoleByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "角色不存在",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取角色成功",
			"data":    role,
		})
	}
}

func updateRoleByIDHandler(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的角色ID",
				"data":    nil,
			})
		}

		role := new(models.Role)
		if err := c.BodyParser(role); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := roleService.UpdateRoleByID(uint(id), role); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "更新角色失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "更新角色成功",
			"data":    role,
		})
	}
}

func deleteRoleByIDHandler(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的角色ID",
				"data":    nil,
			})
		}

		if err := roleService.DeleteRoleByID(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "删除角色失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "删除角色成功",
			"data":    nil,
		})
	}
}
