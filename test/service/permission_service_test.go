package service_test

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionServiceTest(t *testing.T) service.PermissionService {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Permission{})
	assert.NoError(t, err)

	permissionRepo := repository.NewPermissionRepository(db)
	return service.NewPermissionService(permissionRepo)
}

func TestPermissionService_CreatePermission(t *testing.T) {
	permissionService := setupPermissionServiceTest(t)

	permission := &models.Permission{
		Name: "read:users",
	}

	err := permissionService.CreatePermission(permission)
	assert.NoError(t, err)
	assert.NotZero(t, permission.ID)
}

func TestPermissionService_GetPermissions(t *testing.T) {
	permissionService := setupPermissionServiceTest(t)

	// 创建测试权限
	permission1 := &models.Permission{
		Name: "read:users",
	}
	permission2 := &models.Permission{
		Name: "write:users",
	}

	err := permissionService.CreatePermission(permission1)
	assert.NoError(t, err)
	err = permissionService.CreatePermission(permission2)
	assert.NoError(t, err)

	permissions, err := permissionService.GetPermissions()
	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
}

func TestPermissionService_GetPermissionByID(t *testing.T) {
	permissionService := setupPermissionServiceTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := permissionService.CreatePermission(permission)
	assert.NoError(t, err)

	// 查找权限
	foundPermission, err := permissionService.GetPermissionByID(permission.ID)
	assert.NoError(t, err)
	assert.Equal(t, permission.Name, foundPermission.Name)
}

func TestPermissionService_UpdatePermissionByID(t *testing.T) {
	permissionService := setupPermissionServiceTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := permissionService.CreatePermission(permission)
	assert.NoError(t, err)

	// 更新权限
	permission.Name = "write:users"
	err = permissionService.UpdatePermissionByID(permission.ID, permission)
	assert.NoError(t, err)

	// 验证更新
	foundPermission, err := permissionService.GetPermissionByID(permission.ID)
	assert.NoError(t, err)
	assert.Equal(t, "write:users", foundPermission.Name)
}

func TestPermissionService_DeletePermissionByID(t *testing.T) {
	permissionService := setupPermissionServiceTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := permissionService.CreatePermission(permission)
	assert.NoError(t, err)

	// 删除权限
	err = permissionService.DeletePermissionByID(permission.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = permissionService.GetPermissionByID(permission.ID)
	assert.Error(t, err)
}
