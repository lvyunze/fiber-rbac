package repository

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 用户-角色关联测试套件
type UserRoleRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	testUser  *models.User
	testRoles []*models.Role
}

// 初始化测试套件
func (suite *UserRoleRepositoryTestSuite) SetupSuite() {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(suite.T(), err, "打开数据库失败")

	// 自动迁移
	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	assert.NoError(suite.T(), err, "迁移数据库失败")

	suite.db = db
	suite.userRepo = repository.NewUserRepository(db)
	suite.roleRepo = repository.NewRoleRepository(db)
}

// 每个测试前的准备工作
func (suite *UserRoleRepositoryTestSuite) SetupTest() {
	// 清理现有数据
	suite.db.Exec("DELETE FROM user_roles")
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM roles")

	// 创建测试用户
	suite.testUser = &models.User{
		Username: "testuser",
		Password: "password",
	}
	err := suite.userRepo.Create(suite.testUser)
	assert.NoError(suite.T(), err, "创建测试用户失败")

	// 创建测试角色
	suite.testRoles = []*models.Role{
		{Name: "admin"},
		{Name: "editor"},
		{Name: "viewer"},
	}

	for _, role := range suite.testRoles {
		err := suite.roleRepo.Create(role)
		assert.NoError(suite.T(), err, "创建测试角色失败")
	}
}

// 测试分配角色
func (suite *UserRoleRepositoryTestSuite) TestAssignRoles() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles {
		roleIDs = append(roleIDs, role.ID)
	}

	// 分配角色
	err := suite.userRepo.AssignRoles(suite.testUser.ID, roleIDs)
	assert.NoError(suite.T(), err, "分配角色失败")

	// 检查用户的角色
	roles, err := suite.userRepo.GetUserRoles(suite.testUser.ID)
	assert.NoError(suite.T(), err, "获取用户角色失败")
	assert.Equal(suite.T(), len(suite.testRoles), len(roles), "角色数量不匹配")

	// 确认所有角色都被分配
	roleNames := make(map[string]bool)
	for _, role := range roles {
		roleNames[role.Name] = true
	}

	for _, role := range suite.testRoles {
		assert.True(suite.T(), roleNames[role.Name], "角色未被分配: "+role.Name)
	}
}

// 测试移除角色
func (suite *UserRoleRepositoryTestSuite) TestRemoveRoles() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles {
		roleIDs = append(roleIDs, role.ID)
	}

	// 先分配所有角色
	err := suite.userRepo.AssignRoles(suite.testUser.ID, roleIDs)
	assert.NoError(suite.T(), err, "分配角色失败")

	// 只保留第一个角色，移除其他角色
	roleIDsToRemove := roleIDs[1:]
	err = suite.userRepo.RemoveRoles(suite.testUser.ID, roleIDsToRemove)
	assert.NoError(suite.T(), err, "移除角色失败")

	// 检查用户的角色
	roles, err := suite.userRepo.GetUserRoles(suite.testUser.ID)
	assert.NoError(suite.T(), err, "获取用户角色失败")
	assert.Equal(suite.T(), 1, len(roles), "角色数量不匹配")
	assert.Equal(suite.T(), suite.testRoles[0].Name, roles[0].Name, "角色不匹配")
}

// 测试获取用户角色
func (suite *UserRoleRepositoryTestSuite) TestGetUserRoles() {
	// 获取角色ID列表
	var roleIDs []uint
	for _, role := range suite.testRoles[:2] { // 只分配前两个角色
		roleIDs = append(roleIDs, role.ID)
	}

	// 分配角色
	err := suite.userRepo.AssignRoles(suite.testUser.ID, roleIDs)
	assert.NoError(suite.T(), err, "分配角色失败")

	// 获取用户角色
	roles, err := suite.userRepo.GetUserRoles(suite.testUser.ID)
	assert.NoError(suite.T(), err, "获取用户角色失败")
	assert.Equal(suite.T(), 2, len(roles), "角色数量不匹配")

	// 确认正确的角色被分配
	roleNames := make(map[string]bool)
	for _, role := range roles {
		roleNames[role.Name] = true
	}

	for i := 0; i < 2; i++ {
		assert.True(suite.T(), roleNames[suite.testRoles[i].Name], "角色未被分配: "+suite.testRoles[i].Name)
	}
}

// 测试检查用户是否拥有指定角色
func (suite *UserRoleRepositoryTestSuite) TestHasRole() {
	// 为用户分配admin角色
	err := suite.userRepo.AssignRoles(suite.testUser.ID, []uint{suite.testRoles[0].ID})
	assert.NoError(suite.T(), err, "分配角色失败")

	// 检查用户是否拥有admin角色
	hasRole, err := suite.userRepo.HasRole(suite.testUser.ID, "admin")
	assert.NoError(suite.T(), err, "检查角色失败")
	assert.True(suite.T(), hasRole, "用户应该拥有admin角色")

	// 检查用户是否拥有editor角色（未分配）
	hasRole, err = suite.userRepo.HasRole(suite.testUser.ID, "editor")
	assert.NoError(suite.T(), err, "检查角色失败")
	assert.False(suite.T(), hasRole, "用户不应该拥有editor角色")
}

// 运行测试套件
func TestUserRoleRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRoleRepositoryTestSuite))
}
