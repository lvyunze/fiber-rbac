package repository

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 角色-权限关联测试套件
type RolePermissionRepositoryTestSuite struct {
	suite.Suite
	db              *gorm.DB
	roleRepo        repository.RoleRepository
	permissionRepo  repository.PermissionRepository
	testRole        *models.Role
	testPermissions []*models.Permission
}

// 初始化测试套件
func (suite *RolePermissionRepositoryTestSuite) SetupSuite() {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(suite.T(), err, "打开数据库失败")

	// 自动迁移
	err = db.AutoMigrate(&models.Role{}, &models.Permission{})
	assert.NoError(suite.T(), err, "迁移数据库失败")

	suite.db = db
	suite.roleRepo = repository.NewRoleRepository(db)
	suite.permissionRepo = repository.NewPermissionRepository(db)
}

// 每个测试前的准备工作
func (suite *RolePermissionRepositoryTestSuite) SetupTest() {
	// 清理现有数据
	suite.db.Exec("DELETE FROM role_permissions")
	suite.db.Exec("DELETE FROM roles")
	suite.db.Exec("DELETE FROM permissions")

	// 创建测试角色
	suite.testRole = &models.Role{
		Name: "admin",
	}
	err := suite.roleRepo.Create(suite.testRole)
	assert.NoError(suite.T(), err, "创建测试角色失败")

	// 创建测试权限
	suite.testPermissions = []*models.Permission{
		{Name: "user:create"},
		{Name: "user:read"},
		{Name: "user:update"},
		{Name: "user:delete"},
	}

	for _, permission := range suite.testPermissions {
		err := suite.permissionRepo.Create(permission)
		assert.NoError(suite.T(), err, "创建测试权限失败")
	}
}

// 测试分配权限
func (suite *RolePermissionRepositoryTestSuite) TestAssignPermissions() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, permission := range suite.testPermissions {
		permissionIDs = append(permissionIDs, permission.ID)
	}

	// 分配权限
	err := suite.roleRepo.AssignPermissions(suite.testRole.ID, permissionIDs)
	assert.NoError(suite.T(), err, "分配权限失败")

	// 检查角色的权限
	permissions, err := suite.roleRepo.GetRolePermissions(suite.testRole.ID)
	assert.NoError(suite.T(), err, "获取角色权限失败")
	assert.Equal(suite.T(), len(suite.testPermissions), len(permissions), "权限数量不匹配")

	// 确认所有权限都被分配
	permMap := make(map[string]bool)
	for _, perm := range permissions {
		permMap[perm.Name] = true
	}

	for _, perm := range suite.testPermissions {
		assert.True(suite.T(), permMap[perm.Name], "权限未被分配: "+perm.Name)
	}
}

// 测试移除权限
func (suite *RolePermissionRepositoryTestSuite) TestRemovePermissions() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, permission := range suite.testPermissions {
		permissionIDs = append(permissionIDs, permission.ID)
	}

	// 先分配所有权限
	err := suite.roleRepo.AssignPermissions(suite.testRole.ID, permissionIDs)
	assert.NoError(suite.T(), err, "分配权限失败")

	// 只保留前两个权限，移除其他权限
	permissionIDsToRemove := permissionIDs[2:]
	err = suite.roleRepo.RemovePermissions(suite.testRole.ID, permissionIDsToRemove)
	assert.NoError(suite.T(), err, "移除权限失败")

	// 检查角色的权限
	permissions, err := suite.roleRepo.GetRolePermissions(suite.testRole.ID)
	assert.NoError(suite.T(), err, "获取角色权限失败")
	assert.Equal(suite.T(), 2, len(permissions), "权限数量不匹配")

	// 确认正确的权限被保留
	permMap := make(map[string]bool)
	for _, perm := range permissions {
		permMap[perm.Name] = true
	}

	for i := 0; i < 2; i++ {
		perm := suite.testPermissions[i]
		assert.True(suite.T(), permMap[perm.Name], "权限应该被保留: "+perm.Name)
	}
}

// 测试获取角色权限
func (suite *RolePermissionRepositoryTestSuite) TestGetRolePermissions() {
	// 获取特定权限ID
	var permissionIDs []uint
	for i := 0; i < 2; i++ { // 只分配前两个权限
		permissionIDs = append(permissionIDs, suite.testPermissions[i].ID)
	}

	// 分配权限
	err := suite.roleRepo.AssignPermissions(suite.testRole.ID, permissionIDs)
	assert.NoError(suite.T(), err, "分配权限失败")

	// 获取角色权限
	permissions, err := suite.roleRepo.GetRolePermissions(suite.testRole.ID)
	assert.NoError(suite.T(), err, "获取角色权限失败")
	assert.Equal(suite.T(), 2, len(permissions), "权限数量不匹配")

	// 确认正确的权限被分配
	permMap := make(map[string]bool)
	for _, perm := range permissions {
		permMap[perm.Name] = true
	}

	for i := 0; i < 2; i++ {
		perm := suite.testPermissions[i]
		assert.True(suite.T(), permMap[perm.Name], "权限未被分配: "+perm.Name)
	}
}

// 测试检查角色是否拥有指定权限
func (suite *RolePermissionRepositoryTestSuite) TestHasPermission() {
	// 为角色分配创建用户的权限
	createPermission := suite.testPermissions[0] // 用户创建权限
	err := suite.roleRepo.AssignPermissions(suite.testRole.ID, []uint{createPermission.ID})
	assert.NoError(suite.T(), err, "分配权限失败")

	// 检查角色是否有创建用户的权限
	hasPermission, err := suite.roleRepo.HasPermission(suite.testRole.ID, createPermission.Name)
	assert.NoError(suite.T(), err, "检查权限失败")
	assert.True(suite.T(), hasPermission, "角色应该拥有"+createPermission.Name+"权限")

	// 检查角色是否有读取用户的权限（未分配）
	readPermission := suite.testPermissions[1] // 用户读取权限
	hasPermission, err = suite.roleRepo.HasPermission(suite.testRole.ID, readPermission.Name)
	assert.NoError(suite.T(), err, "检查权限失败")
	assert.False(suite.T(), hasPermission, "角色不应该拥有"+readPermission.Name+"权限")
}

// 运行测试套件
func TestRolePermissionRepositorySuite(t *testing.T) {
	suite.Run(t, new(RolePermissionRepositoryTestSuite))
}
