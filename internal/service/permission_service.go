package service

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
)

type PermissionService interface {
	CreatePermission(permission *models.Permission) error
	GetPermissions() ([]models.Permission, error)
	GetPermissionByID(id uint) (*models.Permission, error)
	UpdatePermissionByID(id uint, permission *models.Permission) error
	DeletePermissionByID(id uint) error
}

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) CreatePermission(permission *models.Permission) error {
	return s.repo.Create(permission)
}

func (s *permissionService) GetPermissions() ([]models.Permission, error) {
	return s.repo.FindAll()
}

func (s *permissionService) GetPermissionByID(id uint) (*models.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *permissionService) UpdatePermissionByID(id uint, permission *models.Permission) error {
	return s.repo.Update(id, permission)
}

func (s *permissionService) DeletePermissionByID(id uint) error {
	return s.repo.Delete(id)
}
