package repository

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *models.Role) error
	FindAll() ([]models.Role, error)
	FindByID(id uint) (*models.Role, error)
	Update(id uint, role *models.Role) error
	Delete(id uint) error

	// 角色-权限关联管理
	AssignPermissions(roleID uint, permissionIDs []uint) error
	RemovePermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]models.Permission, error)
	HasPermission(roleID uint, permissionName string) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *roleRepository) Update(id uint, role *models.Role) error {
	return r.db.Model(&models.Role{}).Where("id = ?", id).Updates(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// 分配权限给角色
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 先获取角色
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	// 获取要分配的权限
	var permissions []models.Permission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// 分配权限
	return r.db.Model(&role).Association("Permissions").Append(&permissions)
}

// 移除角色的权限
func (r *roleRepository) RemovePermissions(roleID uint, permissionIDs []uint) error {
	// 先获取角色
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	// 获取要移除的权限
	var permissions []models.Permission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// 移除权限
	return r.db.Model(&role).Association("Permissions").Delete(&permissions)
}

// 获取角色的所有权限
func (r *roleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	var permissions []models.Permission
	if err := r.db.Model(&role).Association("Permissions").Find(&permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// 检查角色是否拥有指定权限
func (r *roleRepository) HasPermission(roleID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, permissionName).
		Count(&count).Error

	return count > 0, err
}
