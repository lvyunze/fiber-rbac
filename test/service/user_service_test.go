package service_test

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserServiceTest(t *testing.T) service.UserService {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	return service.NewUserService(userRepo)
}

func TestUserService_CreateUser(t *testing.T) {
	userService := setupUserServiceTest(t)

	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	err := userService.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserService_GetUsers(t *testing.T) {
	userService := setupUserServiceTest(t)

	// 创建测试用户
	user1 := &models.User{
		Username: "user1",
		Password: "pass1",
	}
	user2 := &models.User{
		Username: "user2",
		Password: "pass2",
	}

	err := userService.CreateUser(user1)
	assert.NoError(t, err)
	err = userService.CreateUser(user2)
	assert.NoError(t, err)

	users, err := userService.GetUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserService_GetUserByID(t *testing.T) {
	userService := setupUserServiceTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := userService.CreateUser(user)
	assert.NoError(t, err)

	// 查找用户
	foundUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Password, foundUser.Password)
}

func TestUserService_UpdateUserByID(t *testing.T) {
	userService := setupUserServiceTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := userService.CreateUser(user)
	assert.NoError(t, err)

	// 更新用户
	user.Username = "updateduser"
	user.Password = "newpassword"
	err = userService.UpdateUserByID(user.ID, user)
	assert.NoError(t, err)

	// 验证更新
	foundUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", foundUser.Username)
	assert.Equal(t, "newpassword", foundUser.Password)
}

func TestUserService_DeleteUserByID(t *testing.T) {
	userService := setupUserServiceTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := userService.CreateUser(user)
	assert.NoError(t, err)

	// 删除用户
	err = userService.DeleteUserByID(user.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = userService.GetUserByID(user.ID)
	assert.Error(t, err)
}
