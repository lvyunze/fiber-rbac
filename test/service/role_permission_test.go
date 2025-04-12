package service

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// 模拟角色仓库
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) FindAll() ([]models.Role, error) {
	args := m.Called()
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByID(id uint) (*models.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) Update(id uint, role *models.Role) error {
	args := m.Called(id, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermissions(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleRepository) HasPermission(roleID uint, permissionName string) (bool, error) {
	args := m.Called(roleID, permissionName)
	return args.Bool(0), args.Error(1)
}

// 角色权限服务测试套件
type RolePermissionServiceTestSuite struct {
	suite.Suite
	mockRoleRepo    *MockRoleRepository
	roleService     service.RoleService
	testRole        *models.Role
	testPermissions []models.Permission
}

// 初始化测试套件
func (suite *RolePermissionServiceTestSuite) SetupTest() {
	suite.mockRoleRepo = new(MockRoleRepository)
	suite.roleService = service.NewRoleService(suite.mockRoleRepo)

	// 创建测试角色
	suite.testRole = &models.Role{
		Model: gorm.Model{ID: 1},
		Name:  "admin",
	}

	// 创建测试权限
	suite.testPermissions = []models.Permission{
		{Model: gorm.Model{ID: 1}, Name: "user:create"},
		{Model: gorm.Model{ID: 2}, Name: "user:read"},
		{Model: gorm.Model{ID: 3}, Name: "user:update"},
		{Model: gorm.Model{ID: 4}, Name: "user:delete"},
	}
}

// 测试分配权限给角色
func (suite *RolePermissionServiceTestSuite) TestAssignPermissionsToRole() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, perm := range suite.testPermissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// 设置模拟预期
	suite.mockRoleRepo.On("AssignPermissions", suite.testRole.ID, permissionIDs).Return(nil)

	// 测试
	err := suite.roleService.AssignPermissionsToRole(suite.testRole.ID, permissionIDs)

	// 断言
	assert.NoError(suite.T(), err, "分配权限应该成功")
	suite.mockRoleRepo.AssertExpectations(suite.T())
}

// 测试移除角色的权限
func (suite *RolePermissionServiceTestSuite) TestRemovePermissionsFromRole() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, perm := range suite.testPermissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// 设置模拟预期
	suite.mockRoleRepo.On("RemovePermissions", suite.testRole.ID, permissionIDs).Return(nil)

	// 测试
	err := suite.roleService.RemovePermissionsFromRole(suite.testRole.ID, permissionIDs)

	// 断言
	assert.NoError(suite.T(), err, "移除权限应该成功")
	suite.mockRoleRepo.AssertExpectations(suite.T())
}

// 测试获取角色权限
func (suite *RolePermissionServiceTestSuite) TestGetRolePermissions() {
	// 设置模拟预期
	suite.mockRoleRepo.On("GetRolePermissions", suite.testRole.ID).Return(suite.testPermissions, nil)

	// 测试
	permissions, err := suite.roleService.GetRolePermissions(suite.testRole.ID)

	// 断言
	assert.NoError(suite.T(), err, "获取角色权限应该成功")
	assert.Equal(suite.T(), len(suite.testPermissions), len(permissions), "权限数量应该匹配")
	suite.mockRoleRepo.AssertExpectations(suite.T())
}

// 测试检查角色是否拥有指定权限
func (suite *RolePermissionServiceTestSuite) TestHasPermission() {
	// 设置模拟预期
	suite.mockRoleRepo.On("HasPermission", suite.testRole.ID, "user:create").Return(true, nil)
	suite.mockRoleRepo.On("HasPermission", suite.testRole.ID, "user:delete").Return(false, nil)

	// 测试拥有的权限
	hasPermission, err := suite.roleService.HasPermission(suite.testRole.ID, "user:create")

	// 断言
	assert.NoError(suite.T(), err, "检查权限应该成功")
	assert.True(suite.T(), hasPermission, "角色应该拥有user:create权限")

	// 测试不拥有的权限
	hasPermission, err = suite.roleService.HasPermission(suite.testRole.ID, "user:delete")

	// 断言
	assert.NoError(suite.T(), err, "检查权限应该成功")
	assert.False(suite.T(), hasPermission, "角色不应该拥有user:delete权限")

	suite.mockRoleRepo.AssertExpectations(suite.T())
}

// 运行测试套件
func TestRolePermissionServiceSuite(t *testing.T) {
	suite.Run(t, new(RolePermissionServiceTestSuite))
}
