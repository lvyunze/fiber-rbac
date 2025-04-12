package repository

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Update(id uint, user *models.User) error
	Delete(id uint) error

	// 用户-角色关联管理
	AssignRoles(userID uint, roleIDs []uint) error
	RemoveRoles(userID uint, roleIDs []uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	HasRole(userID uint, roleName string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(id uint, user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// 分配角色给用户
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// 先获取用户
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 获取要分配的角色
	var roles []models.Role
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	// 分配角色
	return r.db.Model(&user).Association("Roles").Append(&roles)
}

// 移除用户的角色
func (r *userRepository) RemoveRoles(userID uint, roleIDs []uint) error {
	// 先获取用户
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 获取要移除的角色
	var roles []models.Role
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	// 移除角色
	return r.db.Model(&user).Association("Roles").Delete(&roles)
}

// 获取用户的所有角色
func (r *userRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var roles []models.Role
	if err := r.db.Model(&user).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// 检查用户是否拥有指定角色
func (r *userRepository) HasRole(userID uint, roleName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
		Count(&count).Error

	return count > 0, err
}
