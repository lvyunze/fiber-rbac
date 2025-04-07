package repository_test

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoleRepositoryTest(t *testing.T) repository.RoleRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Role{})
	assert.NoError(t, err)

	return repository.NewRoleRepository(db)
}

func TestRoleRepository_Create(t *testing.T) {
	repo := setupRoleRepositoryTest(t)

	role := &models.Role{
		Name: "admin",
	}

	err := repo.Create(role)
	assert.NoError(t, err)
	assert.NotZero(t, role.ID)
}

func TestRoleRepository_FindAll(t *testing.T) {
	repo := setupRoleRepositoryTest(t)

	// 创建测试角色
	role1 := &models.Role{
		Name: "admin",
	}
	role2 := &models.Role{
		Name: "user",
	}

	err := repo.Create(role1)
	assert.NoError(t, err)
	err = repo.Create(role2)
	assert.NoError(t, err)

	roles, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
}

func TestRoleRepository_FindByID(t *testing.T) {
	repo := setupRoleRepositoryTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := repo.Create(role)
	assert.NoError(t, err)

	// 查找角色
	foundRole, err := repo.FindByID(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, role.Name, foundRole.Name)
}

func TestRoleRepository_Update(t *testing.T) {
	repo := setupRoleRepositoryTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := repo.Create(role)
	assert.NoError(t, err)

	// 更新角色
	role.Name = "superadmin"
	err = repo.Update(role.ID, role)
	assert.NoError(t, err)

	// 验证更新
	foundRole, err := repo.FindByID(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, "superadmin", foundRole.Name)
}

func TestRoleRepository_Delete(t *testing.T) {
	repo := setupRoleRepositoryTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := repo.Create(role)
	assert.NoError(t, err)

	// 删除角色
	err = repo.Delete(role.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = repo.FindByID(role.ID)
	assert.Error(t, err)
}
