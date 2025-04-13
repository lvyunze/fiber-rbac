package app

import (
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/handler/auth"
	"github.com/lvyunze/fiber-rbac/internal/handler/permission"
	"github.com/lvyunze/fiber-rbac/internal/handler/role"
	"github.com/lvyunze/fiber-rbac/internal/handler/user"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(app *fiber.App, userService service.UserService, roleService service.RoleService, permissionService service.PermissionService, jwtConfig *config.JWTConfig) {
	// API 版本前缀
	api := app.Group("/api/v1")

	// 认证相关路由 - 无需认证
	authGroup := api.Group("/auth")
	authGroup.Post("/login", auth.NewLoginHandler(userService).Handle)
	authGroup.Post("/refresh", auth.NewRefreshHandler(userService).Handle)

	// 需要认证的路由组
	authRequired := api.Use(middleware.Auth(jwtConfig))

	// 用户个人信息
	authGroup.Get("/profile", middleware.Auth(jwtConfig), auth.NewProfileHandler(userService).Handle)
	authGroup.Post("/check-permission", middleware.Auth(jwtConfig), auth.NewCheckHandler(userService).Handle)

	// 用户管理
	userGroup := authRequired.Group("/users")
	userGroup.Get("/", user.NewListHandler(userService).Handle)
	userGroup.Post("/", user.NewCreateHandler(userService).Handle)
	userGroup.Get("/:id", user.NewDetailHandler(userService).Handle)
	userGroup.Put("/:id", user.NewUpdateHandler(userService).Handle)
	userGroup.Delete("/:id", user.NewDeleteHandler(userService).Handle)
	userGroup.Post("/:id/roles", user.NewAssignRoleHandler(userService).Handle)
	userGroup.Get("/:id/roles", user.NewListRolesHandler(userService).Handle)

	// 角色管理
	roleGroup := authRequired.Group("/roles")
	roleGroup.Get("/", role.NewListHandler(roleService).Handle)
	roleGroup.Post("/", role.NewCreateHandler(roleService).Handle)
	roleGroup.Get("/:id", role.NewDetailHandler(roleService).Handle)
	roleGroup.Put("/:id", role.NewUpdateHandler(roleService).Handle)
	roleGroup.Delete("/:id", role.NewDeleteHandler(roleService).Handle)
	roleGroup.Post("/:id/permissions", role.NewAssignPermissionHandler(roleService).Handle)
	roleGroup.Get("/:id/permissions", role.NewListPermissionsHandler(roleService).Handle)

	// 权限管理
	permissionGroup := authRequired.Group("/permissions")
	permissionGroup.Get("/", permission.NewListHandler(permissionService).Handle)
	permissionGroup.Post("/", permission.NewCreateHandler(permissionService).Handle)
	permissionGroup.Get("/:id", permission.NewDetailHandler(permissionService).Handle)
	permissionGroup.Put("/:id", permission.NewUpdateHandler(permissionService).Handle)
	permissionGroup.Delete("/:id", permission.NewDeleteHandler(permissionService).Handle)
}
