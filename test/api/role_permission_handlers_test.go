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

// 角色权限处理器测试套件
type RolePermissionHandlersTestSuite struct {
	suite.Suite
	app             *fiber.App
	mockRoleService *MockRoleService
	testRole        *models.Role
	testPermissions []models.Permission
}

// 设置测试套件
func (suite *RolePermissionHandlersTestSuite) SetupTest() {
	suite.mockRoleService = new(MockRoleService)

	// 创建Fiber应用
	suite.app = fiber.New()

	// 注册路由
	roles := suite.app.Group("/api/v1/roles")
	v1.RegisterRolePermissionRoutes(roles, suite.mockRoleService)

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
func (suite *RolePermissionHandlersTestSuite) TestAssignPermissionsToRole() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, permission := range suite.testPermissions {
		permissionIDs = append(permissionIDs, permission.ID)
	}

	// 设置模拟预期
	suite.mockRoleService.On("AssignPermissionsToRole", suite.testRole.ID, permissionIDs).Return(nil)

	// 创建请求体
	reqBody := v1.RolePermissionsRequest{
		PermissionIDs: permissionIDs,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 发送请求
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/roles/1/permissions",
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
	assert.Equal(suite.T(), "权限分配成功", response.Message, "响应消息不匹配")

	suite.mockRoleService.AssertExpectations(suite.T())
}

// 测试移除角色的权限
func (suite *RolePermissionHandlersTestSuite) TestRemovePermissionsFromRole() {
	// 获取权限ID列表
	var permissionIDs []uint
	for _, permission := range suite.testPermissions[:2] {
		permissionIDs = append(permissionIDs, permission.ID)
	}

	// 设置模拟预期
	suite.mockRoleService.On("RemovePermissionsFromRole", suite.testRole.ID, permissionIDs).Return(nil)

	// 创建请求体
	reqBody := v1.RolePermissionsRequest{
		PermissionIDs: permissionIDs,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 发送请求
	req, _ := http.NewRequest(
		"DELETE",
		"/api/v1/roles/1/permissions",
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
	assert.Equal(suite.T(), "权限移除成功", response.Message, "响应消息不匹配")

	suite.mockRoleService.AssertExpectations(suite.T())
}

// 测试获取角色权限
func (suite *RolePermissionHandlersTestSuite) TestGetRolePermissions() {
	// 设置模拟预期
	suite.mockRoleService.On("GetRolePermissions", suite.testRole.ID).Return(suite.testPermissions, nil)

	// 发送请求
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/roles/1/permissions",
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
	assert.Equal(suite.T(), "获取角色权限成功", response.Message, "响应消息不匹配")

	// 检查数据
	var permissions []models.Permission
	dataJSON, _ := json.Marshal(response.Data)
	json.Unmarshal(dataJSON, &permissions)

	assert.Equal(suite.T(), len(suite.testPermissions), len(permissions), "权限数量不匹配")

	suite.mockRoleService.AssertExpectations(suite.T())
}

// 测试检查角色是否拥有权限
func (suite *RolePermissionHandlersTestSuite) TestHasPermission() {
	// 设置模拟预期
	suite.mockRoleService.On("HasPermission", suite.testRole.ID, "user:create").Return(true, nil)

	// 发送请求
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/roles/1/has-permission/user:create",
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
	assert.Equal(suite.T(), "检查角色权限成功", response.Message, "响应消息不匹配")

	// 检查数据
	var data map[string]interface{}
	dataJSON, _ := json.Marshal(response.Data)
	json.Unmarshal(dataJSON, &data)

	assert.True(suite.T(), data["has_permission"].(bool), "角色应该拥有user:create权限")

	suite.mockRoleService.AssertExpectations(suite.T())
}

// 运行测试套件
func TestRolePermissionHandlersSuite(t *testing.T) {
	suite.Run(t, new(RolePermissionHandlersTestSuite))
}
