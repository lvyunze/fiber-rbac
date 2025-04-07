package v1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/lvyunze/fiber-rbac/api/v1"
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 自定义错误
var ErrUserNotFound = errors.New("用户不存在")

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) UpdateUserByID(id uint, user *models.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// 测试注册处理程序
func TestRegisterHandler(t *testing.T) {
	// 设置JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 创建模拟用户服务
	mockUserService := new(MockUserService)

	// 创建Fiber应用
	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterAuthRoutes(api, mockUserService)

	tests := []struct {
		name           string
		reqBody        map[string]interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "成功注册",
			reqBody: map[string]interface{}{
				"username": "newuser",
				"password": "password123",
				"email":    "newuser@example.com",
			},
			setupMock: func() {
				mockUserService.On("GetUserByUsername", "newuser").Return(nil, ErrUserNotFound)
				mockUserService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1000), result["code"])
				assert.Equal(t, "success", result["status"])
				assert.Contains(t, result, "data")
				data := result["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name: "用户名已存在",
			reqBody: map[string]interface{}{
				"username": "existinguser",
				"password": "password123",
				"email":    "existing@example.com",
			},
			setupMock: func() {
				existingUser := &models.User{Username: "existinguser"}
				mockUserService.On("GetUserByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1001), result["code"])
				assert.Equal(t, "error", result["status"])
				assert.Equal(t, "用户名已存在", result["message"])
			},
		},
		{
			name: "无效请求",
			reqBody: map[string]interface{}{
				"username": "",
				"password": "short",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1001), result["code"])
				assert.Equal(t, "error", result["status"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟
			tt.setupMock()

			// 创建请求
			reqBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// 执行请求
			resp, _ := app.Test(req)

			// 检查状态和响应
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			tt.checkResponse(t, resp)

			// 清除模拟
			mockUserService.AssertExpectations(t)
		})
	}
}

// 测试登录处理程序
func TestLoginHandler(t *testing.T) {
	// 设置JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 创建模拟用户服务
	mockUserService := new(MockUserService)

	// 创建Fiber应用
	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterAuthRoutes(api, mockUserService)

	// 创建测试用的密码哈希
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		reqBody        map[string]interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name: "成功登录",
			reqBody: map[string]interface{}{
				"username": "validuser",
				"password": "password123",
			},
			setupMock: func() {
				user := &models.User{
					Model:    gorm.Model{ID: 1},
					Username: "validuser",
					Password: string(hashedPassword),
				}
				mockUserService.On("GetUserByUsername", "validuser").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1000), result["code"])
				assert.Equal(t, "success", result["status"])
				assert.Contains(t, result, "data")
				data := result["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name: "用户不存在",
			reqBody: map[string]interface{}{
				"username": "nonexistentuser",
				"password": "password123",
			},
			setupMock: func() {
				mockUserService.On("GetUserByUsername", "nonexistentuser").Return(nil, ErrUserNotFound)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1002), result["code"])
				assert.Equal(t, "error", result["status"])
				assert.Equal(t, "用户名或密码错误", result["message"])
			},
		},
		{
			name: "密码错误",
			reqBody: map[string]interface{}{
				"username": "validuser",
				"password": "wrongpassword",
			},
			setupMock: func() {
				user := &models.User{
					Model:    gorm.Model{ID: 1},
					Username: "validuser",
					Password: string(hashedPassword),
				}
				mockUserService.On("GetUserByUsername", "validuser").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var result map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, float64(1002), result["code"])
				assert.Equal(t, "error", result["status"])
				assert.Equal(t, "用户名或密码错误", result["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟
			tt.setupMock()

			// 创建请求
			reqBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// 执行请求
			resp, _ := app.Test(req)

			// 检查状态和响应
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			tt.checkResponse(t, resp)

			// 清除模拟
			mockUserService.AssertExpectations(t)
		})
	}
}

// 测试令牌刷新处理程序
func TestRefreshTokenHandler(t *testing.T) {
	// 设置JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 创建模拟用户服务
	mockUserService := new(MockUserService)

	// 创建Fiber应用
	app := fiber.New()
	api := app.Group("/api/v1")
	v1.RegisterAuthRoutes(api, mockUserService)

	// 生成有效令牌
	validUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "validuser",
	}
	validToken, _ := utils.GenerateToken(validUser.ID, validUser.Username)

	// 创建一个无效令牌（错误格式）
	invalidToken := "invalid.token.string"

	// 创建一个空令牌
	emptyToken := ""

	tests := []struct {
		name            string
		reqBody         map[string]interface{}
		expectedStatus  int
		expectedCode    float64
		expectedResult  string
		expectedMessage string
	}{
		{
			name: "成功刷新令牌",
			reqBody: map[string]interface{}{
				"token": validToken,
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    1000,
			expectedResult:  "success",
			expectedMessage: "",
		},
		{
			name: "无效令牌",
			reqBody: map[string]interface{}{
				"token": invalidToken,
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    1008,
			expectedResult:  "error",
			expectedMessage: "无效的令牌",
		},
		{
			name: "未提供令牌",
			reqBody: map[string]interface{}{
				"token": emptyToken,
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    1006,
			expectedResult:  "error",
			expectedMessage: "未提供令牌",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			reqBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// 执行请求
			resp, _ := app.Test(req)

			// 检查状态码
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// 解析响应
			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)

			// 检查状态和错误代码
			assert.Equal(t, tt.expectedCode, result["code"])
			assert.Equal(t, tt.expectedResult, result["status"])

			// 如果期望成功
			if tt.expectedResult == "success" {
				// 验证数据结构
				assert.Contains(t, result, "data")
				data := result["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.NotEmpty(t, data["token"])
			} else {
				// 验证错误消息
				assert.Equal(t, tt.expectedMessage, result["message"])
			}
		})
	}
}
