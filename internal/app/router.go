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

	// 用户个人信息和权限检查
	authGroup.Post("/profile", middleware.Auth(jwtConfig), auth.NewProfileHandler(userService).Handle)
	authGroup.Post("/check-permission", middleware.Auth(jwtConfig), auth.NewCheckHandler(userService).Handle)

	// 用户管理
	userGroup := authRequired.Group("/users")
	userGroup.Post("/list", user.NewListHandler(userService).Handle)
	userGroup.Post("/create", user.NewCreateHandler(userService).Handle)
	userGroup.Post("/detail", user.NewDetailHandler(userService).Handle)
	userGroup.Post("/update", user.NewUpdateHandler(userService).Handle)
	userGroup.Post("/delete", user.NewDeleteHandler(userService).Handle)
	userGroup.Post("/assign-roles", user.NewAssignRoleHandler(userService).Handle)
	userGroup.Post("/list-roles", user.NewListRolesHandler(userService).Handle)

	// 角色管理
	roleGroup := authRequired.Group("/roles")
	roleGroup.Post("/list", role.NewListHandler(roleService).Handle)
	roleGroup.Post("/create", role.NewCreateHandler(roleService).Handle)
	roleGroup.Post("/detail", role.NewDetailHandler(roleService).Handle)
	roleGroup.Post("/update", role.NewUpdateHandler(roleService).Handle)
	roleGroup.Post("/delete", role.NewDeleteHandler(roleService).Handle)
	roleGroup.Post("/assign-permissions", role.NewAssignPermissionHandler(roleService).Handle)
	roleGroup.Post("/list-permissions", role.NewListPermissionsHandler(roleService).Handle)

	// 权限管理
	permissionGroup := authRequired.Group("/permissions")
	permissionGroup.Post("/list", permission.NewListHandler(permissionService).Handle)
	permissionGroup.Post("/create", permission.NewCreateHandler(permissionService).Handle)
	permissionGroup.Post("/detail", permission.NewDetailHandler(permissionService).Handle)
	permissionGroup.Post("/update", permission.NewUpdateHandler(permissionService).Handle)
	permissionGroup.Post("/delete", permission.NewDeleteHandler(permissionService).Handle)
}
