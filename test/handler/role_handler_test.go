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

func setupRoleHandlerTest(t *testing.T) (*fiber.App, service.RoleService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Role{})
	assert.NoError(t, err)

	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)

	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterRoleRoutes(api, roleService)

	return app, roleService
}

func TestCreateRoleHandler(t *testing.T) {
	app, _ := setupRoleHandlerTest(t)

	// 创建角色请求
	role := models.Role{
		Name: "管理员",
	}
	body, _ := json.Marshal(role)
	req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
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
	assert.Equal(t, "角色创建成功", response.Message)
	assert.NotNil(t, response.Data)
}

func TestGetRolesHandler(t *testing.T) {
	app, roleService := setupRoleHandlerTest(t)

	// 创建测试角色
	role1 := &models.Role{
		Name: "管理员",
	}
	role2 := &models.Role{
		Name: "用户",
	}
	roleService.CreateRole(role1)
	roleService.CreateRole(role2)

	// 获取角色列表请求
	req := httptest.NewRequest("GET", "/api/v1/roles", nil)

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
	assert.Equal(t, "获取角色列表成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证角色数量
	var roles []models.Role
	rolesData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(rolesData, &roles)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
}

func TestGetRoleByIDHandler(t *testing.T) {
	app, roleService := setupRoleHandlerTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "管理员",
	}
	roleService.CreateRole(role)

	// 获取角色请求
	req := httptest.NewRequest("GET", "/api/v1/roles/1", nil)

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
	assert.Equal(t, "获取角色成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证角色数据
	var foundRole models.Role
	roleData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(roleData, &foundRole)
	assert.NoError(t, err)
	assert.Equal(t, "管理员", foundRole.Name)
}

func TestGetRoleByIDHandler_NotFound(t *testing.T) {
	app, _ := setupRoleHandlerTest(t)

	// 获取不存在的角色请求
	req := httptest.NewRequest("GET", "/api/v1/roles/999", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeRoleNotFound, response.Code)
	assert.Equal(t, "角色不存在", response.Message)
	assert.Nil(t, response.Data)
}

func TestUpdateRoleByIDHandler(t *testing.T) {
	app, roleService := setupRoleHandlerTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "管理员",
	}
	roleService.CreateRole(role)

	// 更新角色请求
	updatedRole := models.Role{
		Name: "超级管理员",
	}
	body, _ := json.Marshal(updatedRole)
	req := httptest.NewRequest("PUT", "/api/v1/roles/1", bytes.NewBuffer(body))
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
	assert.Equal(t, "更新角色成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证更新后的角色
	var foundRole models.Role
	roleData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(roleData, &foundRole)
	assert.NoError(t, err)
	assert.Equal(t, "超级管理员", foundRole.Name)
}

func TestDeleteRoleByIDHandler(t *testing.T) {
	app, roleService := setupRoleHandlerTest(t)

	// 创建测试角色
	role := &models.Role{
		Name: "管理员",
	}
	roleService.CreateRole(role)

	// 删除角色请求
	req := httptest.NewRequest("DELETE", "/api/v1/roles/1", nil)

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
	assert.Equal(t, "删除角色成功", response.Message)
	assert.Nil(t, response.Data)

	// 验证角色已删除
	_, err = roleService.GetRoleByID(1)
	assert.Error(t, err)
}
