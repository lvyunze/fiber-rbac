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

func setupRoleServiceTest(t *testing.T) service.RoleService {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Role{})
	assert.NoError(t, err)

	roleRepo := repository.NewRoleRepository(db)
	return service.NewRoleService(roleRepo)
}

func TestRoleService_CreateRole(t *testing.T) {
	roleService := setupRoleServiceTest(t)

	role := &models.Role{
		Name: "admin",
	}

	err := roleService.CreateRole(role)
	assert.NoError(t, err)
	assert.NotZero(t, role.ID)
}

func TestRoleService_GetRoles(t *testing.T) {
	roleService := setupRoleServiceTest(t)

	// 创建测试角色
	role1 := &models.Role{
		Name: "admin",
	}
	role2 := &models.Role{
		Name: "user",
	}

	err := roleService.CreateRole(role1)
	assert.NoError(t, err)
	err = roleService.CreateRole(role2)
	assert.NoError(t, err)

	roles, err := roleService.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
}

func TestRoleService_GetRoleByID(t *testing.T) {
	roleService := setupRoleServiceTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := roleService.CreateRole(role)
	assert.NoError(t, err)

	// 查找角色
	foundRole, err := roleService.GetRoleByID(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, role.Name, foundRole.Name)
}

func TestRoleService_UpdateRoleByID(t *testing.T) {
	roleService := setupRoleServiceTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := roleService.CreateRole(role)
	assert.NoError(t, err)

	// 更新角色
	role.Name = "superadmin"
	err = roleService.UpdateRoleByID(role.ID, role)
	assert.NoError(t, err)

	// 验证更新
	foundRole, err := roleService.GetRoleByID(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, "superadmin", foundRole.Name)
}

func TestRoleService_DeleteRoleByID(t *testing.T) {
	roleService := setupRoleServiceTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "admin",
	}
	err := roleService.CreateRole(role)
	assert.NoError(t, err)

	// 删除角色
	err = roleService.DeleteRoleByID(role.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = roleService.GetRoleByID(role.ID)
	assert.Error(t, err)
}
