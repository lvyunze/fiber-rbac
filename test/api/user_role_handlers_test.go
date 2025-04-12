package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

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

// 用户角色处理器测试套件
type UserRoleHandlersTestSuite struct {
	suite.Suite
	app             *fiber.App
	mockUserService *MockUserService
	testUser        *models.User
	testRoles       []models.Role
}

// 设置测试套件
func (suite *UserRoleHandlersTestSuite) SetupTest() {
	suite.mockUserService = new(MockUserService)

	// 创建Fiber应用
	suite.app = fiber.New()

	// 注册路由
	users := suite.app.Group("/api/v1/users")
	v1.RegisterUserRoleRoutes(users, suite.mockUserService)

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
func (suite *UserRoleHandlersTestSuite) TestAssignRolesToUser() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles {
		roleIDs = append(roleIDs, role.ID)
	}

	// 设置模拟预期
	suite.mockUserService.On("AssignRolesToUser", suite.testUser.ID, roleIDs).Return(nil)

	// 创建请求体
	reqBody := v1.UserRolesRequest{
		RoleIDs: roleIDs,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 发送请求
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/users/1/roles",
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 解析响应
	var response utils.Response
	json.NewDecoder(resp.Body).Decode(&response)

	// 验证响应
	assert.Equal(suite.T(), utils.StatusSuccess, response.Status, "响应状态应该是成功")
	assert.Equal(suite.T(), "角色分配成功", response.Message, "响应消息不匹配")

	suite.mockUserService.AssertExpectations(suite.T())
}

// 测试移除用户的角色
func (suite *UserRoleHandlersTestSuite) TestRemoveRolesFromUser() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles[:2] {
		roleIDs = append(roleIDs, role.ID)
	}

	// 设置模拟预期
	suite.mockUserService.On("RemoveRolesFromUser", suite.testUser.ID, roleIDs).Return(nil)

	// 创建请求体
	reqBody := v1.UserRolesRequest{
		RoleIDs: roleIDs,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 发送请求
	req, _ := http.NewRequest(
		"DELETE",
		"/api/v1/users/1/roles",
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 解析响应
	var response utils.Response
	json.NewDecoder(resp.Body).Decode(&response)

	// 验证响应
	assert.Equal(suite.T(), utils.StatusSuccess, response.Status, "响应状态应该是成功")
	assert.Equal(suite.T(), "角色移除成功", response.Message, "响应消息不匹配")

	suite.mockUserService.AssertExpectations(suite.T())
}

// 测试获取用户角色
func (suite *UserRoleHandlersTestSuite) TestGetUserRoles() {
	// 设置模拟预期
	suite.mockUserService.On("GetUserRoles", suite.testUser.ID).Return(suite.testRoles, nil)

	// 发送请求
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/users/1/roles",
		nil,
	)

	// 执行请求
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 解析响应
	var response utils.Response
	json.NewDecoder(resp.Body).Decode(&response)

	// 验证响应
	assert.Equal(suite.T(), utils.StatusSuccess, response.Status, "响应状态应该是成功")
	assert.Equal(suite.T(), "获取用户角色成功", response.Message, "响应消息不匹配")

	// 检查数据
	var roles []models.Role
	dataJSON, _ := json.Marshal(response.Data)
	json.Unmarshal(dataJSON, &roles)

	assert.Equal(suite.T(), len(suite.testRoles), len(roles), "角色数量不匹配")

	suite.mockUserService.AssertExpectations(suite.T())
}

// 测试检查用户是否拥有角色
func (suite *UserRoleHandlersTestSuite) TestHasRole() {
	// 设置模拟预期
	suite.mockUserService.On("HasRole", suite.testUser.ID, "admin").Return(true, nil)

	// 发送请求
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/users/1/has-role/admin",
		nil,
	)

	// 执行请求
	resp, err := suite.app.Test(req)

	// 断言
	assert.NoError(suite.T(), err, "请求应该成功")
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode, "状态码应该是200")

	// 解析响应
	var response utils.Response
	json.NewDecoder(resp.Body).Decode(&response)

	// 验证响应
	assert.Equal(suite.T(), utils.StatusSuccess, response.Status, "响应状态应该是成功")
	assert.Equal(suite.T(), "检查用户角色成功", response.Message, "响应消息不匹配")

	// 检查数据
	var data map[string]interface{}
	dataJSON, _ := json.Marshal(response.Data)
	json.Unmarshal(dataJSON, &data)

	assert.True(suite.T(), data["has_role"].(bool), "用户应该拥有admin角色")

	suite.mockUserService.AssertExpectations(suite.T())
}

// 运行测试套件
func TestUserRoleHandlersSuite(t *testing.T) {
	suite.Run(t, new(UserRoleHandlersTestSuite))
}
