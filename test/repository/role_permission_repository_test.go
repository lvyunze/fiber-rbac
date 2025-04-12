package repository_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RolePermissionRepositoryTestSuite struct {
	suite.Suite
	db              *gorm.DB
	roleRepo        repository.RoleRepository
	permissionRepo  repository.PermissionRepository
	testRole        *models.Role
	testPermissions []*models.Permission
	testIndex       int
	rand            *rand.Rand
}

func (s *RolePermissionRepositoryTestSuite) SetupSuite() {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	s.rand = rand.New(rand.NewSource(seed))
}

func (s *RolePermissionRepositoryTestSuite) SetupTest() {
	var err error
	// 为每个测试用例使用独立的内存数据库
	dbName := fmt.Sprintf("file::memory:rptest%d?mode=memory&cache=shared", s.rand.Int())
	s.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	s.Require().NoError(err)

	// 自动迁移所有相关模型
	err = s.db.AutoMigrate(&models.Role{}, &models.Permission{})
	s.Require().NoError(err)

	// 初始化仓库
	s.roleRepo = repository.NewRoleRepository(s.db)
	s.permissionRepo = repository.NewPermissionRepository(s.db)

	s.testIndex++
	// 创建测试角色和权限
	s.testRole = s.createTestRole()
	s.testPermissions = s.createTestPermissions()
}

func (s *RolePermissionRepositoryTestSuite) randomString(prefix string) string {
	return prefix + "_" + strconv.Itoa(s.testIndex) + "_" + strconv.FormatInt(s.rand.Int63(), 16)
}

func (s *RolePermissionRepositoryTestSuite) createTestRole() *models.Role {
	role := &models.Role{
		Name: s.randomString("admin"),
	}
	err := s.roleRepo.Create(role)
	s.Require().NoError(err)
	return role
}

func (s *RolePermissionRepositoryTestSuite) createTestPermissions() []*models.Permission {
	permissions := []*models.Permission{
		{Name: s.randomString("create")},
		{Name: s.randomString("read")},
		{Name: s.randomString("update")},
	}

	for i := range permissions {
		err := s.permissionRepo.Create(permissions[i])
		s.Require().NoError(err)
	}
	return permissions
}

func (s *RolePermissionRepositoryTestSuite) TestAssignPermissions() {
	// 分配第一个和第二个权限给角色
	permissionIDs := []uint{s.testPermissions[0].ID, s.testPermissions[1].ID}
	err := s.roleRepo.AssignPermissions(s.testRole.ID, permissionIDs)
	s.Require().NoError(err)

	// 验证权限是否被正确分配
	rolePermissions, err := s.roleRepo.GetRolePermissions(s.testRole.ID)
	s.Require().NoError(err)
	s.Equal(2, len(rolePermissions))

	// 验证权限ID是否正确
	permissionIDMap := make(map[uint]bool)
	for _, permission := range rolePermissions {
		permissionIDMap[permission.ID] = true
	}
	s.True(permissionIDMap[s.testPermissions[0].ID])
	s.True(permissionIDMap[s.testPermissions[1].ID])
}

func (s *RolePermissionRepositoryTestSuite) TestRemovePermissions() {
	// 首先分配所有权限
	permissionIDs := []uint{s.testPermissions[0].ID, s.testPermissions[1].ID, s.testPermissions[2].ID}
	err := s.roleRepo.AssignPermissions(s.testRole.ID, permissionIDs)
	s.Require().NoError(err)

	// 移除第一个权限
	permissionIDsToRemove := []uint{s.testPermissions[0].ID}
	err = s.roleRepo.RemovePermissions(s.testRole.ID, permissionIDsToRemove)
	s.Require().NoError(err)

	// 验证权限是否被正确移除
	rolePermissions, err := s.roleRepo.GetRolePermissions(s.testRole.ID)
	s.Require().NoError(err)
	s.Equal(2, len(rolePermissions))

	// 验证剩余的权限ID是否正确
	permissionIDMap := make(map[uint]bool)
	for _, permission := range rolePermissions {
		permissionIDMap[permission.ID] = true
	}
	s.False(permissionIDMap[s.testPermissions[0].ID])
	s.True(permissionIDMap[s.testPermissions[1].ID])
	s.True(permissionIDMap[s.testPermissions[2].ID])
}

func (s *RolePermissionRepositoryTestSuite) TestGetRolePermissions() {
	// 分配第二个和第三个权限给角色
	permissionIDs := []uint{s.testPermissions[1].ID, s.testPermissions[2].ID}
	err := s.roleRepo.AssignPermissions(s.testRole.ID, permissionIDs)
	s.Require().NoError(err)

	// 获取角色权限
	rolePermissions, err := s.roleRepo.GetRolePermissions(s.testRole.ID)
	s.Require().NoError(err)
	s.Equal(2, len(rolePermissions))

	// 验证权限名称是否正确
	permissionNameMap := make(map[string]bool)
	for _, permission := range rolePermissions {
		permissionNameMap[permission.Name] = true
	}
	s.True(permissionNameMap[s.testPermissions[1].Name])
	s.True(permissionNameMap[s.testPermissions[2].Name])
}

func (s *RolePermissionRepositoryTestSuite) TestHasPermission() {
	// 分配第一个权限给角色
	permissionIDs := []uint{s.testPermissions[0].ID}
	err := s.roleRepo.AssignPermissions(s.testRole.ID, permissionIDs)
	s.Require().NoError(err)

	// 验证角色是否拥有第一个权限
	hasPermission, err := s.roleRepo.HasPermission(s.testRole.ID, s.testPermissions[0].Name)
	s.Require().NoError(err)
	s.True(hasPermission)

	// 验证角色是否没有第二个权限
	hasPermission, err = s.roleRepo.HasPermission(s.testRole.ID, s.testPermissions[1].Name)
	s.Require().NoError(err)
	s.False(hasPermission)
}

func TestRolePermissionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RolePermissionRepositoryTestSuite))
}
