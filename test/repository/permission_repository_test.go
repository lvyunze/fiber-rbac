package repository_test

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionRepositoryTest(t *testing.T) repository.PermissionRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Permission{})
	assert.NoError(t, err)

	return repository.NewPermissionRepository(db)
}

func TestPermissionRepository_Create(t *testing.T) {
	repo := setupPermissionRepositoryTest(t)

	permission := &models.Permission{
		Name: "read:users",
	}

	err := repo.Create(permission)
	assert.NoError(t, err)
	assert.NotZero(t, permission.ID)
}

func TestPermissionRepository_FindAll(t *testing.T) {
	repo := setupPermissionRepositoryTest(t)

	// 创建测试权限
	permission1 := &models.Permission{
		Name: "read:users",
	}
	permission2 := &models.Permission{
		Name: "write:users",
	}

	err := repo.Create(permission1)
	assert.NoError(t, err)
	err = repo.Create(permission2)
	assert.NoError(t, err)

	permissions, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
}

func TestPermissionRepository_FindByID(t *testing.T) {
	repo := setupPermissionRepositoryTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := repo.Create(permission)
	assert.NoError(t, err)

	// 查找权限
	foundPermission, err := repo.FindByID(permission.ID)
	assert.NoError(t, err)
	assert.Equal(t, permission.Name, foundPermission.Name)
}

func TestPermissionRepository_Update(t *testing.T) {
	repo := setupPermissionRepositoryTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := repo.Create(permission)
	assert.NoError(t, err)

	// 更新权限
	permission.Name = "write:users"
	err = repo.Update(permission.ID, permission)
	assert.NoError(t, err)

	// 验证更新
	foundPermission, err := repo.FindByID(permission.ID)
	assert.NoError(t, err)
	assert.Equal(t, "write:users", foundPermission.Name)
}

func TestPermissionRepository_Delete(t *testing.T) {
	repo := setupPermissionRepositoryTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "read:users",
	}
	err := repo.Create(permission)
	assert.NoError(t, err)

	// 删除权限
	err = repo.Delete(permission.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = repo.FindByID(permission.ID)
	assert.Error(t, err)
}
