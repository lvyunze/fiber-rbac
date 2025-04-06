package main

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/config"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	// Load configuration
	_ = config.LoadConfig()

	// Initialize database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	userService := service.NewUserService(*userRepo)
	roleService := service.NewRoleService(*roleRepo)
	permissionService := service.NewPermissionService(*permissionRepo)

	// Register routes
	router := app.Group("/api/v1")
	v1.RegisterUserRoutes(router, userService)
	v1.RegisterRoleRoutes(router, roleService)
	v1.RegisterPermissionRoutes(router, permissionService)

	// Start server
	app.Listen(":3000")
}
