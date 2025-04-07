package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionHandlerTest(t *testing.T) (*fiber.App, service.PermissionService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Permission{})
	assert.NoError(t, err)

	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo)

	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterPermissionRoutes(api, permissionService)

	return app, permissionService
}

func TestCreatePermissionHandler(t *testing.T) {
	app, _ := setupPermissionHandlerTest(t)

	// 创建权限请求
	permission := models.Permission{
		Name: "用户管理",
	}
	body, _ := json.Marshal(permission)
	req := httptest.NewRequest("POST", "/api/v1/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusSuccess, response.Status)
	assert.Equal(t, utils.CodeSuccess, response.Code)
	assert.Equal(t, "权限创建成功", response.Message)
	assert.NotNil(t, response.Data)
}

func TestGetPermissionsHandler(t *testing.T) {
	app, permissionService := setupPermissionHandlerTest(t)

	// 创建测试权限
	permission1 := &models.Permission{
		Name: "用户管理",
	}
	permission2 := &models.Permission{
		Name: "角色管理",
	}
	permissionService.CreatePermission(permission1)
	permissionService.CreatePermission(permission2)

	// 获取权限列表请求
	req := httptest.NewRequest("GET", "/api/v1/permissions", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusSuccess, response.Status)
	assert.Equal(t, utils.CodeSuccess, response.Code)
	assert.Equal(t, "获取权限列表成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证权限数量
	var permissions []models.Permission
	permissionsData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(permissionsData, &permissions)
	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
}

func TestGetPermissionByIDHandler(t *testing.T) {
	app, permissionService := setupPermissionHandlerTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "用户管理",
	}
	permissionService.CreatePermission(permission)

	// 获取权限请求
	req := httptest.NewRequest("GET", "/api/v1/permissions/1", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusSuccess, response.Status)
	assert.Equal(t, utils.CodeSuccess, response.Code)
	assert.Equal(t, "获取权限成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证权限数据
	var foundPermission models.Permission
	permissionData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(permissionData, &foundPermission)
	assert.NoError(t, err)
	assert.Equal(t, "用户管理", foundPermission.Name)
}

func TestGetPermissionByIDHandler_NotFound(t *testing.T) {
	app, _ := setupPermissionHandlerTest(t)

	// 获取不存在的权限请求
	req := httptest.NewRequest("GET", "/api/v1/permissions/999", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodePermissionNotFound, response.Code)
	assert.Equal(t, "权限不存在", response.Message)
	assert.Nil(t, response.Data)
}

func TestUpdatePermissionByIDHandler(t *testing.T) {
	app, permissionService := setupPermissionHandlerTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "用户管理",
	}
	permissionService.CreatePermission(permission)

	// 更新权限请求
	updatedPermission := models.Permission{
		Name: "用户管理高级",
	}
	body, _ := json.Marshal(updatedPermission)
	req := httptest.NewRequest("PUT", "/api/v1/permissions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusSuccess, response.Status)
	assert.Equal(t, utils.CodeSuccess, response.Code)
	assert.Equal(t, "更新权限成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证更新后的权限
	var foundPermission models.Permission
	permissionData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(permissionData, &foundPermission)
	assert.NoError(t, err)
	assert.Equal(t, "用户管理高级", foundPermission.Name)
}

func TestDeletePermissionByIDHandler(t *testing.T) {
	app, permissionService := setupPermissionHandlerTest(t)

	// 创建测试权限
	permission := &models.Permission{
		Name: "用户管理",
	}
	permissionService.CreatePermission(permission)

	// 删除权限请求
	req := httptest.NewRequest("DELETE", "/api/v1/permissions/1", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusSuccess, response.Status)
	assert.Equal(t, utils.CodeSuccess, response.Code)
	assert.Equal(t, "删除权限成功", response.Message)
	assert.Nil(t, response.Data)

	// 验证权限已删除
	_, err = permissionService.GetPermissionByID(1)
	assert.Error(t, err)
}
