package service_test

import (
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 测试用户服务创建用户功能
func TestUserService_Create(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name           string
		request        *schema.CreateUserRequest
		mockSetup      func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository)
		expectedID     uint64
		expectedError  error
		expectedCalled bool
	}{
		{
			name: "创建用户成功",
			request: &schema.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// 模拟用户不存在
				mockUserRepo.On("GetByUsername", "testuser").Return(nil, nil)
				
				// 模拟邮箱不存在
				mockUserRepo.On("GetByEmail", "test@example.com").Return(nil, nil)
				
				// 模拟创建用户成功
				mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
					user := args.Get(0).(*model.User)
					user.ID = 1 // 设置ID
				}).Return(nil)
			},
			expectedID:     1,
			expectedError:  nil,
			expectedCalled: true,
		},
		{
			name: "用户已存在",
			request: &schema.CreateUserRequest{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// 模拟用户已存在
				existingUser := &model.User{ID: 1, Username: "existinguser"}
				mockUserRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedID:     0,
			expectedError:  errors.ErrUserExists,
			expectedCalled: false,
		},
		{
			name: "邮箱已存在",
			request: &schema.CreateUserRequest{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// 模拟用户不存在
				mockUserRepo.On("GetByUsername", "testuser").Return(nil, nil)
				
				// 模拟邮箱已存在
				existingUser := &model.User{ID: 1, Email: "existing@example.com"}
				mockUserRepo.On("GetByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedID:     0,
			expectedError:  errors.ErrEmailExists,
			expectedCalled: false,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓库
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
			
			// 设置模拟行为
			tt.mockSetup(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo)
			
			// 创建JWT配置
			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}
			
			// 创建用户服务
			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)
			
			// 调用创建用户方法
			id, err := userService.Create(tt.request)
			
			// 断言结果
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)
			
			// 断言模拟仓库的Create方法是否被调用
			if tt.expectedCalled {
				mockUserRepo.AssertCalled(t, "Create", mock.AnythingOfType("*model.User"))
			} else {
				mockUserRepo.AssertNotCalled(t, "Create", mock.AnythingOfType("*model.User"))
			}
		})
	}
}

// 测试用户服务登录功能
func TestUserService_Login(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name          string
		request       *schema.LoginRequest
		mockSetup     func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository)
		expectedError error
		expectedToken bool
	}{
		{
			name: "登录成功",
			request: &schema.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				existingUser := &model.User{
					ID:       1,
					Username: "testuser",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$dDmrbhFKvY/rYmKkxsiDNw$h0QDgvpBVhD79Uk7C0LEa3Jr3pVJ4v3vaqUFmPlY+Xg", // "password123"
				}
				mockUserRepo.On("GetByUsername", "testuser").Return(existingUser, nil)
				// refresh token mock
				mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*model.UserRefreshToken")).Return(nil)
			},
			expectedError: nil,
			expectedToken: true,
		},
		{
			name: "用户不存在",
			request: &schema.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetByUsername", "nonexistent").Return(nil, nil)
			},
			expectedError: errors.ErrInvalidCredentials,
			expectedToken: false,
		},
		{
			name: "密码错误",
			request: &schema.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				existingUser := &model.User{
					ID:       1,
					Username: "testuser",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$dDmrbhFKvY/rYmKkxsiDNw$h0QDgvpBVhD79Uk7C0LEa3Jr3pVJ4v3vaqUFmPlY+Xg", // "password123"
				}
				mockUserRepo.On("GetByUsername", "testuser").Return(existingUser, nil)
			},
			expectedError: errors.ErrInvalidCredentials,
			expectedToken: false,
		},
	}

	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)

			// 设置模拟行为
			tt.mockSetup(mockUserRepo, mockRefreshTokenRepo)

			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}

			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)

			response, err := userService.Login(tt.request)

			assert.Equal(t, tt.expectedError, err)

			if tt.expectedToken {
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, jwtConfig.Expire, response.ExpiresIn)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

// 测试用户服务获取用户信息功能
func TestUserService_GetProfile(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name          string
		userID        uint64
		mockSetup     func(mockUserRepo *mocks.MockUserRepository)
		expectedError error
		expectedUser  bool
	}{
		{
			name:   "获取用户信息成功",
			userID: 1,
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				existingUser := &model.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					CreatedAt: 1617235200,
					Roles:     []model.Role{},
				}
				mockUserRepo.On("GetByID", uint64(1)).Return(existingUser, nil)
			},
			expectedError: nil,
			expectedUser:  true,
		},
		{
			name:   "用户不存在",
			userID: 999,
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetByID", uint64(999)).Return(nil, nil)
			},
			expectedError: errors.ErrUserNotFound,
			expectedUser:  false,
		},
	}

	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)

			// 设置模拟行为
			tt.mockSetup(mockUserRepo)

			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}

			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)

			user, err := userService.GetProfile(tt.userID)

			assert.Equal(t, tt.expectedError, err)

			if tt.expectedUser {
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

// TestUserService_List 分页列表测试
func TestUserService_List(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)
	mockPermRepo := new(mocks.MockPermissionRepository)
	mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
	service := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, &config.JWTConfig{})

	tests := []struct {
		name       string
		page       int
		pageSize   int
		keyword    string
		mockSetup  func()
		expectResp *schema.ListUserResponse
		expectErr  error
	}{
		{
			name:     "正常分页返回",
			page:     1,
			pageSize: 2,
			keyword:  "",
			mockSetup: func() {
				users := []*model.User{{ID: 1, Username: "u1"}, {ID: 2, Username: "u2"}}
				mockUserRepo.On("List", 1, 2, "").Return(users, int64(2), nil)
			},
			expectResp: &schema.ListUserResponse{
				Total:      2,
				Page:       1,
				PageSize:   2,
				TotalPages: 1,
				Items: []schema.UserResponse{{ID: 1, Username: "u1", Roles: []schema.RoleSimple{}}, {ID: 2, Username: "u2", Roles: []schema.RoleSimple{}}},
			},
			expectErr: nil,
		},
		{
			name:     "无数据",
			page:     1,
			pageSize: 10,
			keyword:  "",
			mockSetup: func() {
				mockUserRepo.On("List", 1, 10, "").Return([]*model.User{}, int64(0), nil)
			},
			expectResp: &schema.ListUserResponse{
				Total:      0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
				Items:      []schema.UserResponse{},
			},
			expectErr: nil,
		},
		{
			name:     "数据库异常",
			page:     1,
			pageSize: 10,
			keyword:  "",
			mockSetup: func() {
				mockUserRepo.On("List", 1, 10, "").Return([]*model.User(nil), int64(0), errors.ErrDB)
			},
			expectResp: nil,
			expectErr: errors.ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil // 清理历史
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			resp, err := service.List(&schema.ListUserRequest{Page: tt.page, PageSize: tt.pageSize, Keyword: tt.keyword})
			assert.Equal(t, tt.expectErr, err)
			assert.Equal(t, tt.expectResp, resp)
		})
	}
}
