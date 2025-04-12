package service

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserByID(id uint, user *models.User) error
	DeleteUserByID(id uint) error

	// 用户-角色关联管理
	AssignRolesToUser(userID uint, roleIDs []uint) error
	RemoveRolesFromUser(userID uint, roleIDs []uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	HasRole(userID uint, roleName string) (bool, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.repo.Create(user)
}

func (s *userService) GetUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *userService) UpdateUserByID(id uint, user *models.User) error {
	return s.repo.Update(id, user)
}

func (s *userService) DeleteUserByID(id uint) error {
	return s.repo.Delete(id)
}

// 分配角色给用户
func (s *userService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	return s.repo.AssignRoles(userID, roleIDs)
}

// 移除用户的角色
func (s *userService) RemoveRolesFromUser(userID uint, roleIDs []uint) error {
	return s.repo.RemoveRoles(userID, roleIDs)
}

// 获取用户的所有角色
func (s *userService) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.repo.GetUserRoles(userID)
}

// 检查用户是否拥有指定角色
func (s *userService) HasRole(userID uint, roleName string) (bool, error) {
	return s.repo.HasRole(userID, roleName)
}
