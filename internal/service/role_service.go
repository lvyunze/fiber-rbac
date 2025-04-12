package service

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
)

type RoleService interface {
	CreateRole(role *models.Role) error
	GetRoles() ([]models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
	UpdateRoleByID(id uint, role *models.Role) error
	DeleteRoleByID(id uint) error

	// 角色-权限关联管理
	AssignPermissionsToRole(roleID uint, permissionIDs []uint) error
	RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]models.Permission, error)
	HasPermission(roleID uint, permissionName string) (bool, error)
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(role *models.Role) error {
	return s.repo.Create(role)
}

func (s *roleService) GetRoles() ([]models.Role, error) {
	return s.repo.FindAll()
}

func (s *roleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleService) UpdateRoleByID(id uint, role *models.Role) error {
	return s.repo.Update(id, role)
}

func (s *roleService) DeleteRoleByID(id uint) error {
	return s.repo.Delete(id)
}

// 分配权限给角色
func (s *roleService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	return s.repo.AssignPermissions(roleID, permissionIDs)
}

// 移除角色的权限
func (s *roleService) RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error {
	return s.repo.RemovePermissions(roleID, permissionIDs)
}

// 获取角色的所有权限
func (s *roleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	return s.repo.GetRolePermissions(roleID)
}

// 检查角色是否拥有指定权限
func (s *roleService) HasPermission(roleID uint, permissionName string) (bool, error) {
	return s.repo.HasPermission(roleID, permissionName)
}
