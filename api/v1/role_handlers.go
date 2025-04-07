package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// RegisterRoleRoutes 注册角色相关的路由
func RegisterRoleRoutes(router fiber.Router, roleService service.RoleService) {
	roles := router.Group("/roles")

	roles.Post("/", createRoleHandler(roleService))
	roles.Get("/", getRolesHandler(roleService))
	roles.Get("/:id", getRoleByIDHandler(roleService))
	roles.Put("/:id", updateRoleByIDHandler(roleService))
	roles.Delete("/:id", deleteRoleByIDHandler(roleService))
}

func createRoleHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := new(models.Role)
		if err := c.BodyParser(role); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := roleService.CreateRole(role); err != nil {
			return utils.ServerError(c, "创建角色失败")
		}

		return utils.SuccessResponse(c, "角色创建成功", role)
	}
}

func getRolesHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := roleService.GetRoles()
		if err != nil {
			return utils.ServerError(c, "获取角色列表失败")
		}

		return utils.SuccessResponse(c, "获取角色列表成功", roles)
	}
}

func getRoleByIDHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		role, err := roleService.GetRoleByID(uint(id))
		if err != nil {
			return utils.NotFoundError(c, "角色不存在", utils.ErrRoleNotFound)
		}

		return utils.SuccessResponse(c, "获取角色成功", role)
	}
}

func updateRoleByIDHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		role := new(models.Role)
		if err := c.BodyParser(role); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := roleService.UpdateRoleByID(uint(id), role); err != nil {
			return utils.ServerError(c, "更新角色失败")
		}

		return utils.SuccessResponse(c, "更新角色成功", role)
	}
}

func deleteRoleByIDHandler(roleService service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的角色ID")
		}

		if err := roleService.DeleteRoleByID(uint(id)); err != nil {
			return utils.ServerError(c, "删除角色失败")
		}

		return utils.SuccessResponse(c, "删除角色成功", nil)
	}
}
