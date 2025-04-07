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

func setupUserHandlerTest(t *testing.T) (*fiber.App, service.UserService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterUserRoutes(api, userService)

	return app, userService
}

func TestCreateUserHandler(t *testing.T) {
	app, _ := setupUserHandlerTest(t)

	// 创建用户请求
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
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
	assert.Equal(t, "用户创建成功", response.Message)
	assert.NotNil(t, response.Data)
}

func TestGetUsersHandler(t *testing.T) {
	app, userService := setupUserHandlerTest(t)

	// 创建测试用户
	user1 := &models.User{
		Username: "user1",
		Password: "pass1",
	}
	user2 := &models.User{
		Username: "user2",
		Password: "pass2",
	}
	userService.CreateUser(user1)
	userService.CreateUser(user2)

	// 获取用户列表请求
	req := httptest.NewRequest("GET", "/api/v1/users", nil)

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
	assert.Equal(t, "获取用户列表成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证用户数量
	var users []models.User
	usersData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(usersData, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetUserByIDHandler(t *testing.T) {
	app, userService := setupUserHandlerTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	userService.CreateUser(user)

	// 获取用户请求
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)

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
	assert.Equal(t, "获取用户成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证用户数据
	var foundUser models.User
	userData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(userData, &foundUser)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", foundUser.Username)
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	app, _ := setupUserHandlerTest(t)

	// 获取不存在的用户请求
	req := httptest.NewRequest("GET", "/api/v1/users/999", nil)

	// 发送请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 解析响应
	var response utils.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeUserNotFound, response.Code)
	assert.Equal(t, "用户不存在", response.Message)
	assert.Nil(t, response.Data)
}

func TestUpdateUserByIDHandler(t *testing.T) {
	app, userService := setupUserHandlerTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	userService.CreateUser(user)

	// 更新用户请求
	updatedUser := models.User{
		Username: "updateduser",
		Password: "newpassword",
	}
	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(body))
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
	assert.Equal(t, "更新用户成功", response.Message)
	assert.NotNil(t, response.Data)

	// 验证更新后的用户
	var foundUser models.User
	userData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(userData, &foundUser)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", foundUser.Username)
}

func TestDeleteUserByIDHandler(t *testing.T) {
	app, userService := setupUserHandlerTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	userService.CreateUser(user)

	// 删除用户请求
	req := httptest.NewRequest("DELETE", "/api/v1/users/1", nil)

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
	assert.Equal(t, "删除用户成功", response.Message)
	assert.Nil(t, response.Data)

	// 验证用户已删除
	_, err = userService.GetUserByID(1)
	assert.Error(t, err)
}
