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

type UserRoleRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	testUser  *models.User
	testRoles []*models.Role
	testIndex int
	rand      *rand.Rand
}

func (s *UserRoleRepositoryTestSuite) SetupSuite() {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	s.rand = rand.New(rand.NewSource(seed))
}

func (s *UserRoleRepositoryTestSuite) SetupTest() {
	var err error
	// 为每个测试用例使用独立的内存数据库
	dbName := fmt.Sprintf("file::memory:urtest%d?mode=memory&cache=shared", s.rand.Int())
	s.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	s.Require().NoError(err)

	// 自动迁移所有相关模型
	err = s.db.AutoMigrate(&models.User{}, &models.Role{})
	s.Require().NoError(err)

	// 初始化仓库
	s.userRepo = repository.NewUserRepository(s.db)
	s.roleRepo = repository.NewRoleRepository(s.db)

	s.testIndex++
	// 创建测试用户和角色
	s.testUser = s.createTestUser()
	s.testRoles = s.createTestRoles()
}

func (s *UserRoleRepositoryTestSuite) randomString(prefix string) string {
	return prefix + "_" + strconv.Itoa(s.testIndex) + "_" + strconv.FormatInt(s.rand.Int63(), 16)
}

func (s *UserRoleRepositoryTestSuite) createTestUser() *models.User {
	user := &models.User{
		Username: s.randomString("user"),
		Password: "hashedpassword",
	}
	err := s.userRepo.Create(user)
	s.Require().NoError(err)
	return user
}

func (s *UserRoleRepositoryTestSuite) createTestRoles() []*models.Role {
	roles := []*models.Role{
		{Name: s.randomString("admin")},
		{Name: s.randomString("editor")},
		{Name: s.randomString("viewer")},
	}

	for i := range roles {
		err := s.roleRepo.Create(roles[i])
		s.Require().NoError(err)
	}
	return roles
}

func (s *UserRoleRepositoryTestSuite) TestAssignRoles() {
	// 分配第一个和第二个角色给用户
	roleIDs := []uint{s.testRoles[0].ID, s.testRoles[1].ID}
	err := s.userRepo.AssignRoles(s.testUser.ID, roleIDs)
	s.Require().NoError(err)

	// 验证角色是否被正确分配
	userRoles, err := s.userRepo.GetUserRoles(s.testUser.ID)
	s.Require().NoError(err)
	s.Equal(2, len(userRoles))

	// 验证角色ID是否正确
	roleIDMap := make(map[uint]bool)
	for _, role := range userRoles {
		roleIDMap[role.ID] = true
	}
	s.True(roleIDMap[s.testRoles[0].ID])
	s.True(roleIDMap[s.testRoles[1].ID])
}

func (s *UserRoleRepositoryTestSuite) TestRemoveRoles() {
	// 首先分配所有角色
	roleIDs := []uint{s.testRoles[0].ID, s.testRoles[1].ID, s.testRoles[2].ID}
	err := s.userRepo.AssignRoles(s.testUser.ID, roleIDs)
	s.Require().NoError(err)

	// 移除第一个角色
	roleIDsToRemove := []uint{s.testRoles[0].ID}
	err = s.userRepo.RemoveRoles(s.testUser.ID, roleIDsToRemove)
	s.Require().NoError(err)

	// 验证角色是否被正确移除
	userRoles, err := s.userRepo.GetUserRoles(s.testUser.ID)
	s.Require().NoError(err)
	s.Equal(2, len(userRoles))

	// 验证剩余的角色ID是否正确
	roleIDMap := make(map[uint]bool)
	for _, role := range userRoles {
		roleIDMap[role.ID] = true
	}
	s.False(roleIDMap[s.testRoles[0].ID])
	s.True(roleIDMap[s.testRoles[1].ID])
	s.True(roleIDMap[s.testRoles[2].ID])
}

func (s *UserRoleRepositoryTestSuite) TestGetUserRoles() {
	// 分配第二个和第三个角色给用户
	roleIDs := []uint{s.testRoles[1].ID, s.testRoles[2].ID}
	err := s.userRepo.AssignRoles(s.testUser.ID, roleIDs)
	s.Require().NoError(err)

	// 获取用户角色
	userRoles, err := s.userRepo.GetUserRoles(s.testUser.ID)
	s.Require().NoError(err)
	s.Equal(2, len(userRoles))

	// 验证角色名称是否正确
	roleNameMap := make(map[string]bool)
	for _, role := range userRoles {
		roleNameMap[role.Name] = true
	}
	s.True(roleNameMap[s.testRoles[1].Name])
	s.True(roleNameMap[s.testRoles[2].Name])
}

func (s *UserRoleRepositoryTestSuite) TestHasRole() {
	// 分配第一个角色给用户
	roleIDs := []uint{s.testRoles[0].ID}
	err := s.userRepo.AssignRoles(s.testUser.ID, roleIDs)
	s.Require().NoError(err)

	// 验证用户是否拥有第一个角色
	hasRole, err := s.userRepo.HasRole(s.testUser.ID, s.testRoles[0].Name)
	s.Require().NoError(err)
	s.True(hasRole)

	// 验证用户是否没有第二个角色
	hasRole, err = s.userRepo.HasRole(s.testUser.ID, s.testRoles[1].Name)
	s.Require().NoError(err)
	s.False(hasRole)
}

func TestUserRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleRepositoryTestSuite))
}
