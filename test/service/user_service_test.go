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

// u6d4bu8bd5u7528u6237u521bu5efau529fu80fd
func TestUserService_Create(t *testing.T) {
	// u6d4bu8bd5u7528u4f8bu8868
	tests := []struct {
		name           string
		request        *schema.CreateUserRequest
		mockSetup      func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository)
		expectedID     uint64
		expectedError  error
		expectedCalled bool
	}{
		{
			name: "u6210u529fu521bu5efau7528u6237",
			request: &schema.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// u6a21u62dfu7528u6237u540du4e0du5b58u5728
				mockUserRepo.On("GetByUsername", "testuser").Return(nil, nil)
				
				// u6a21u62dfu90aeu7bb1u4e0du5b58u5728
				mockUserRepo.On("GetByEmail", "test@example.com").Return(nil, nil)
				
				// u6a21u62dfu521bu5efau7528u6237u6210u529f
				mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
					user := args.Get(0).(*model.User)
					user.ID = 1 // u8bbeu7f6eID
				}).Return(nil)
			},
			expectedID:     1,
			expectedError:  nil,
			expectedCalled: true,
		},
		{
			name: "u7528u6237u540du5df2u5b58u5728",
			request: &schema.CreateUserRequest{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// u6a21u62dfu7528u6237u540du5df2u5b58u5728
				existingUser := &model.User{ID: 1, Username: "existinguser"}
				mockUserRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedID:     0,
			expectedError:  errors.ErrUserExists,
			expectedCalled: false,
		},
		{
			name: "u90aeu7bb1u5df2u5b58u5728",
			request: &schema.CreateUserRequest{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository, mockRoleRepo *mocks.MockRoleRepository, mockPermRepo *mocks.MockPermissionRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// u6a21u62dfu7528u6237u540du4e0du5b58u5728
				mockUserRepo.On("GetByUsername", "testuser").Return(nil, nil)
				
				// u6a21u62dfu90aeu7bb1u5df2u5b58u5728
				existingUser := &model.User{ID: 1, Email: "existing@example.com"}
				mockUserRepo.On("GetByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedID:     0,
			expectedError:  errors.ErrEmailExists,
			expectedCalled: false,
		},
	}

	// u6267u884cu6d4bu8bd5u7528u4f8b
	for _, tt := range tests {
		tt := tt // u9632u6b62u95edu5305u95eeu9898
		t.Run(tt.name, func(t *testing.T) {
			// u521bu5efau6a21u62dfu5bf9u8c61
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
			
			// u8bbeu7f6eu6a21u62dfu884cu4e3a
			tt.mockSetup(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo)
			
			// u521bu5efaJWTu914du7f6e
			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}
			
			// u521bu5efau670du52a1u5b9eu4f8b
			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)
			
			// u8c03u7528u88abu6d4bu8bd5u7684u65b9u6cd5
			id, err := userService.Create(tt.request)
			
			// u9a8cu8bc1u7ed3u679c
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)
			
			// u9a8cu8bc1u6a21u62dfu5bf9u8c61u7684u65b9u6cd5u662fu5426u88abu8c03u7528
			if tt.expectedCalled {
				mockUserRepo.AssertCalled(t, "Create", mock.AnythingOfType("*model.User"))
			} else {
				mockUserRepo.AssertNotCalled(t, "Create", mock.AnythingOfType("*model.User"))
			}
		})
	}
}

// u6d4bu8bd5u7528u6237u767bu5f55u529fu80fd
func TestUserService_Login(t *testing.T) {
	// u6d4bu8bd5u7528u4f8bu8868
	tests := []struct {
		name          string
		request       *schema.LoginRequest
		mockSetup     func(mockUserRepo *mocks.MockUserRepository)
		expectedError error
		expectedToken bool
	}{
		{
			name: "u767bu5f55u6210u529f",
			request: &schema.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				// u6a21u62dfu7528u6237u5b58u5728
				// u6ce8u610fuff1au8fd9u91ccu4f7fu7528u4e86u5df2u7ecfu54c8u5e0cu8fc7u7684u5bc6u7801uff0cu5bf9u5e94u660eu6587u662f"password123"
				existingUser := &model.User{
					ID:       1,
					Username: "testuser",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$dDmrbhFKvY/rYmKkxsiDNw$h0QDgvpBVhD79Uk7C0LEa3Jr3pVJ4v3vaqUFmPlY+Xg", // u5bf9u5e94"password123"
				}
				mockUserRepo.On("GetByUsername", "testuser").Return(existingUser, nil)
			},
			expectedError: nil,
			expectedToken: true,
		},
		{
			name: "u7528u6237u4e0du5b58u5728",
			request: &schema.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				// u6a21u62dfu7528u6237u4e0du5b58u5728
				mockUserRepo.On("GetByUsername", "nonexistent").Return(nil, nil)
			},
			expectedError: errors.ErrInvalidCredentials,
			expectedToken: false,
		},
		{
			name: "u5bc6u7801u9519u8bef",
			request: &schema.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				// u6a21u62dfu7528u6237u5b58u5728u4f46u5bc6u7801u9519u8bef
				existingUser := &model.User{
					ID:       1,
					Username: "testuser",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$dDmrbhFKvY/rYmKkxsiDNw$h0QDgvpBVhD79Uk7C0LEa3Jr3pVJ4v3vaqUFmPlY+Xg", // u5bf9u5e94"password123"
				}
				mockUserRepo.On("GetByUsername", "testuser").Return(existingUser, nil)
			},
			expectedError: errors.ErrInvalidCredentials,
			expectedToken: false,
		},
	}

	// u6267u884cu6d4bu8bd5u7528u4f8b
	for _, tt := range tests {
		tt := tt // u9632u6b62u95edu5305u95eeu9898
		t.Run(tt.name, func(t *testing.T) {
			// u521bu5efau6a21u62dfu5bf9u8c61
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
			
			// u8bbeu7f6eu6a21u62dfu884cu4e3a
			tt.mockSetup(mockUserRepo)
			
			// u521bu5efaJWTu914du7f6e
			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}
			
			// u521bu5efau670du52a1u5b9eu4f8b
			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)
			
			// u8c03u7528u88abu6d4bu8bd5u7684u65b9u6cd5
			response, err := userService.Login(tt.request)
			
			// u9a8cu8bc1u7ed3u679c
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

// u6d4bu8bd5u83b7u53d6u7528u6237u4e2au4ebau4fe1u606fu529fu80fd
func TestUserService_GetProfile(t *testing.T) {
	// u6d4bu8bd5u7528u4f8bu8868
	tests := []struct {
		name          string
		userID        uint64
		mockSetup     func(mockUserRepo *mocks.MockUserRepository)
		expectedError error
		expectedUser  bool
	}{
		{
			name:   "u6210u529fu83b7u53d6u7528u6237u4fe1u606f",
			userID: 1,
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				// u6a21u62dfu7528u6237u5b58u5728
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
			name:   "u7528u6237u4e0du5b58u5728",
			userID: 999,
			mockSetup: func(mockUserRepo *mocks.MockUserRepository) {
				// u6a21u62dfu7528u6237u4e0du5b58u5728
				mockUserRepo.On("GetByID", uint64(999)).Return(nil, nil)
			},
			expectedError: errors.ErrUserNotFound,
			expectedUser:  false,
		},
	}

	// u6267u884cu6d4bu8bd5u7528u4f8b
	for _, tt := range tests {
		tt := tt // u9632u6b62u95edu5305u95eeu9898
		t.Run(tt.name, func(t *testing.T) {
			// u521bu5efau6a21u62dfu5bf9u8c61
			mockUserRepo := new(mocks.MockUserRepository)
			mockRoleRepo := new(mocks.MockRoleRepository)
			mockPermRepo := new(mocks.MockPermissionRepository)
			mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
			
			// u8bbeu7f6eu6a21u62dfu884cu4e3a
			tt.mockSetup(mockUserRepo)
			
			// u521bu5efaJWTu914du7f6e
			jwtConfig := &config.JWTConfig{
				Secret: "test-secret",
				Expire: 3600,
			}
			
			// u521bu5efau670du52a1u5b9eu4f8b
			userService := service.NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockRefreshTokenRepo, jwtConfig)
			
			// u8c03u7528u88abu6d4bu8bd5u7684u65b9u6cd5
			user, err := userService.GetProfile(tt.userID)
			
			// u9a8cu8bc1u7ed3u679c
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
