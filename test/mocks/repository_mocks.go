package mocks

import (
	"github.com/lvyunze/fiber-rbac/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository 用户仓库的模拟实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint64) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) List(page, pageSize int, keyword string) ([]*model.User, int64, error) {
	args := m.Called(page, pageSize, keyword)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) AddRoles(userID uint64, roleIDs []uint64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveRoles(userID uint64, roleIDs []uint64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRoles(userID uint64, roleIDs []uint64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserWithRoles(userID uint64) (*model.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// MockRoleRepository 角色仓库的模拟实现
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(role *model.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(role *model.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(id uint64) (*model.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByCode(code string) (*model.Role, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(name string) (*model.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) List(page, pageSize int, keyword string) ([]*model.Role, int64, error) {
	args := m.Called(page, pageSize, keyword)
	return args.Get(0).([]*model.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) AddPermissions(roleID uint64, permissionIDs []uint64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermissions(roleID uint64, permissionIDs []uint64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) UpdatePermissions(roleID uint64, permissionIDs []uint64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRoleWithPermissions(roleID uint64) (*model.Role, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetUsersByRoleID(roleID uint64) ([]*model.User, error) {
	args := m.Called(roleID)
	return args.Get(0).([]*model.User), args.Error(1)
}

// MockRefreshTokenRepository 刷新令牌仓库mock
//go:generate mockery --name=RefreshTokenRepository --output=. --outpkg=mocks --case=underscore
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(refreshToken *model.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Update(refreshToken *model.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByID(id uint64) (*model.RefreshToken, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) GetByToken(token string) (*model.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) List(page, pageSize int, keyword string) ([]*model.RefreshToken, int64, error) {
	args := m.Called(page, pageSize, keyword)
	return args.Get(0).([]*model.RefreshToken), args.Get(1).(int64), args.Error(2)
}

// MockPermissionRepository 权限仓库的模拟实现
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(permission *model.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Update(permission *model.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByID(id uint64) (*model.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByCode(code string) (*model.Permission, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByName(name string) (*model.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) List(page, pageSize int, keyword string) ([]*model.Permission, int64, error) {
	args := m.Called(page, pageSize, keyword)
	return args.Get(0).([]*model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) GetRolesByPermissionID(permissionID uint64) ([]*model.Role, error) {
	args := m.Called(permissionID)
	return args.Get(0).([]*model.Role), args.Error(1)
}

func (m *MockPermissionRepository) GetByNames(names []string) ([]*model.Permission, error) {
	args := m.Called(names)
	return args.Get(0).([]*model.Permission), args.Error(1)
}
