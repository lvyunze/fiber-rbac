package repository_test

import (
	"testing"

	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserRepositoryTest(t *testing.T) repository.UserRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	return repository.NewUserRepository(db)
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupUserRepositoryTest(t)

	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindAll(t *testing.T) {
	repo := setupUserRepositoryTest(t)

	// 创建测试用户
	user1 := &models.User{
		Username: "user1",
		Password: "pass1",
	}
	user2 := &models.User{
		Username: "user2",
		Password: "pass2",
	}

	err := repo.Create(user1)
	assert.NoError(t, err)
	err = repo.Create(user2)
	assert.NoError(t, err)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := setupUserRepositoryTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// 查找用户
	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Password, foundUser.Password)
}

func TestUserRepository_Update(t *testing.T) {
	repo := setupUserRepositoryTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// 更新用户
	user.Username = "updateduser"
	user.Password = "newpassword"
	err = repo.Update(user.ID, user)
	assert.NoError(t, err)

	// 验证更新
	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", foundUser.Username)
	assert.Equal(t, "newpassword", foundUser.Password)
}

func TestUserRepository_Delete(t *testing.T) {
	repo := setupUserRepositoryTest(t)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// 删除用户
	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	// 验证删除
	_, err = repo.FindByID(user.ID)
	assert.Error(t, err)
}
