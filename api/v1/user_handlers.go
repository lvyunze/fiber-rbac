package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// RegisterUserRoutes 注册用户相关的路由
func RegisterUserRoutes(router fiber.Router, userService service.UserService) {
	users := router.Group("/users")

	users.Post("/", createUserHandler(userService))
	users.Get("/", getUsersHandler(userService))
	users.Get("/:id", getUserByIDHandler(userService))
	users.Put("/:id", updateUserByIDHandler(userService))
	users.Delete("/:id", deleteUserByIDHandler(userService))
}

func createUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := userService.CreateUser(user); err != nil {
			return utils.ServerError(c, "创建用户失败")
		}

		return utils.SuccessResponse(c, "用户创建成功", user)
	}
}

func getUsersHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userService.GetUsers()
		if err != nil {
			return utils.ServerError(c, "获取用户列表失败")
		}

		return utils.SuccessResponse(c, "获取用户列表成功", users)
	}
}

func getUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		user, err := userService.GetUserByID(uint(id))
		if err != nil {
			return utils.NotFoundError(c, "用户不存在", utils.ErrUserNotFound)
		}

		return utils.SuccessResponse(c, "获取用户成功", user)
	}
}

func updateUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return utils.BadRequestError(c, "无法解析JSON")
		}

		if err := userService.UpdateUserByID(uint(id), user); err != nil {
			return utils.ServerError(c, "更新用户失败")
		}

		return utils.SuccessResponse(c, "更新用户成功", user)
	}
}

func deleteUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.BadRequestError(c, "无效的用户ID")
		}

		if err := userService.DeleteUserByID(uint(id)); err != nil {
			return utils.ServerError(c, "删除用户失败")
		}

		return utils.SuccessResponse(c, "删除用户成功", nil)
	}
}
