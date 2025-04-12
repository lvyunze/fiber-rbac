package middleware

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// 模拟角色服务
type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) CreateRole(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleService) GetRoles() ([]models.Role, error) {
	args := m.Called()
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleService) GetRoleByID(id uint) (*models.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleService) UpdateRoleByID(id uint, role *models.Role) error {
	args := m.Called(id, role)
	return args.Error(0)
}

func (m *MockRoleService) DeleteRoleByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleService) RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleService) HasPermission(roleID uint, permissionName string) (bool, error) {
	args := m.Called(roleID, permissionName)
	return args.Bool(0), args.Error(1)
}

// 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUserByID(id uint, user *models.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserService) RemoveRolesFromUser(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserService) GetUserRoles(userID uint) ([]models.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockUserService) HasRole(userID uint, roleName string) (bool, error) {
	args := m.Called(userID, roleName)
	return args.Bool(0), args.Error(1)
}

// 权限中间件测试套件
type PermissionMiddlewareTestSuite struct {
	suite.Suite
	app             *fiber.App
	mockUserService *MockUserService
	mockRoleService *MockRoleService
	testUser        *models.User
	testRoles       []models.Role
}

// 设置测试套件
func (suite *PermissionMiddlewareTestSuite) SetupTest() {
	suite.mockUserService = new(MockUserService)
	suite.mockRoleService = new(MockRoleService)

	// 创建Fiber应用
	suite.app = fiber.New()

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
	}
}

// 测试需要特定角色的中间件
func (suite *PermissionMiddlewareTestSuite) TestRequireRole() {
	// 设置路由
	suite.app.Use(func(c *fiber.Ctx) error {
		// 模拟JWT中间件设置的用户ID
		c.Locals("userID", uint(1))
		return c.Next()
	})

	// 需要admin角色
	suite.app.Get("/admin-only", middleware.RequireRole(suite.mockUserService, "admin"), func(c *fiber.Ctx) error {
		return c.SendString("Admin area")
	})

	// 设置模拟预期 - 有admin角色的情况
	suite.mockUserService.On("HasRole", uint(1), "admin").Return(true, nil).Once()

	// 发送请求 - 应该成功
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 设置模拟预期 - 没有admin角色的情况
	suite.mockUserService.On("HasRole", uint(1), "admin").Return(false, nil).Once()

	// 发送请求 - 应该失败
	req, _ = http.NewRequest("GET", "/admin-only", nil)
	resp, err = suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusForbidden, resp.StatusCode, "状态码应该是403")

	suite.mockUserService.AssertExpectations(suite.T())
}

// 测试需要特定权限的中间件
func (suite *PermissionMiddlewareTestSuite) TestRequirePermission() {
	// 设置路由
	suite.app.Use(func(c *fiber.Ctx) error {
		// 模拟JWT中间件设置的用户ID
		c.Locals("userID", uint(1))
		return c.Next()
	})

	// 需要user:create权限
	suite.app.Post("/users", middleware.RequirePermission(suite.mockUserService, suite.mockRoleService, "user:create"), func(c *fiber.Ctx) error {
		return c.SendString("User created")
	})

	// 设置模拟预期 - 有admin角色且该角色有user:create权限的情况
	suite.mockUserService.On("GetUserRoles", uint(1)).Return([]models.Role{{Model: gorm.Model{ID: 1}, Name: "admin"}}, nil).Once()
	suite.mockRoleService.On("HasPermission", uint(1), "user:create").Return(true, nil).Once()

	// 发送请求 - 应该成功
	req, _ := http.NewRequest("POST", "/users", nil)
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 设置模拟预期 - 没有user:create权限的情况
	suite.mockUserService.On("GetUserRoles", uint(1)).Return([]models.Role{{Model: gorm.Model{ID: 2}, Name: "editor"}}, nil).Once()
	suite.mockRoleService.On("HasPermission", uint(2), "user:create").Return(false, nil).Once()

	// 发送请求 - 应该失败
	req, _ = http.NewRequest("POST", "/users", nil)
	resp, err = suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusForbidden, resp.StatusCode, "状态码应该是403")

	suite.mockUserService.AssertExpectations(suite.T())
	suite.mockRoleService.AssertExpectations(suite.T())
}

// 运行测试套件
func TestPermissionMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(PermissionMiddlewareTestSuite))
}
