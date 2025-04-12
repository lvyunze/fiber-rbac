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

// 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(id uint, user *models.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveRoles(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockUserRepository) HasRole(userID uint, roleName string) (bool, error) {
	args := m.Called(userID, roleName)
	return args.Bool(0), args.Error(1)
}

// 用户角色服务测试套件
type UserRoleServiceTestSuite struct {
	suite.Suite
	mockUserRepo *MockUserRepository
	userService  service.UserService
	testUser     *models.User
	testRoles    []models.Role
}

// 初始化测试套件
func (suite *UserRoleServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.userService = service.NewUserService(suite.mockUserRepo)

	// 创建测试用户
	suite.testUser = &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Password: "password",
	}

	// 创建测试角色
	suite.testRoles = []models.Role{
		{Model: gorm.Model{ID: 1}, Name: "admin"},
		{Model: gorm.Model{ID: 2}, Name: "editor"},
		{Model: gorm.Model{ID: 3}, Name: "viewer"},
	}
}

// 测试分配角色给用户
func (suite *UserRoleServiceTestSuite) TestAssignRolesToUser() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles {
		roleIDs = append(roleIDs, role.ID)
	}

	// 设置模拟预期
	suite.mockUserRepo.On("AssignRoles", suite.testUser.ID, roleIDs).Return(nil)

	// 测试
	err := suite.userService.AssignRolesToUser(suite.testUser.ID, roleIDs)

	// 断言
	assert.NoError(suite.T(), err, "分配角色应该成功")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// 测试移除用户的角色
func (suite *UserRoleServiceTestSuite) TestRemoveRolesFromUser() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles {
		roleIDs = append(roleIDs, role.ID)
	}

	// 设置模拟预期
	suite.mockUserRepo.On("RemoveRoles", suite.testUser.ID, roleIDs).Return(nil)

	// 测试
	err := suite.userService.RemoveRolesFromUser(suite.testUser.ID, roleIDs)

	// 断言
	assert.NoError(suite.T(), err, "移除角色应该成功")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// 测试获取用户角色
func (suite *UserRoleServiceTestSuite) TestGetUserRoles() {
	// 设置模拟预期
	suite.mockUserRepo.On("GetUserRoles", suite.testUser.ID).Return(suite.testRoles, nil)

	// 测试
	roles, err := suite.userService.GetUserRoles(suite.testUser.ID)

	// 断言
	assert.NoError(suite.T(), err, "获取用户角色应该成功")
	assert.Equal(suite.T(), len(suite.testRoles), len(roles), "角色数量应该匹配")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// 测试检查用户是否拥有角色
func (suite *UserRoleServiceTestSuite) TestHasRole() {
	// 设置模拟预期
	suite.mockUserRepo.On("HasRole", suite.testUser.ID, "admin").Return(true, nil)
	suite.mockUserRepo.On("HasRole", suite.testUser.ID, "manager").Return(false, nil)

	// 测试拥有的角色
	hasRole, err := suite.userService.HasRole(suite.testUser.ID, "admin")

	// 断言
	assert.NoError(suite.T(), err, "检查角色应该成功")
	assert.True(suite.T(), hasRole, "用户应该拥有admin角色")

	// 测试不拥有的角色
	hasRole, err = suite.userService.HasRole(suite.testUser.ID, "manager")

	// 断言
	assert.NoError(suite.T(), err, "检查角色应该成功")
	assert.False(suite.T(), hasRole, "用户不应该拥有manager角色")

	suite.mockUserRepo.AssertExpectations(suite.T())
}

// 运行测试套件
func TestUserRoleServiceSuite(t *testing.T) {
	suite.Run(t, new(UserRoleServiceTestSuite))
}
