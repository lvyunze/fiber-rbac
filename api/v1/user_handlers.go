package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := userService.CreateUser(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "创建用户失败",
				"data":    nil,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "用户创建成功",
			"data":    user,
		})
	}
}

func getUsersHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userService.GetUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "获取用户列表失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取用户列表成功",
			"data":    users,
		})
	}
}

func getUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的用户ID",
				"data":    nil,
			})
		}

		user, err := userService.GetUserByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "用户不存在",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "获取用户成功",
			"data":    user,
		})
	}
}

func updateUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的用户ID",
				"data":    nil,
			})
		}

		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无法解析JSON",
				"data":    nil,
			})
		}

		if err := userService.UpdateUserByID(uint(id), user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "更新用户失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "更新用户成功",
			"data":    user,
		})
	}
}

func deleteUserByIDHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "无效的用户ID",
				"data":    nil,
			})
		}

		if err := userService.DeleteUserByID(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "删除用户失败",
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "删除用户成功",
			"data":    nil,
		})
	}
}
