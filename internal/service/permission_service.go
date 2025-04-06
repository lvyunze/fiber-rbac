package service

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
)

type PermissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) CreatePermission(permission *models.Permission) error {
	return s.repo.Create(permission)
}

func (s *PermissionService) GetPermissions() ([]models.Permission, error) {
	return s.repo.FindAll()
}

func (s *PermissionService) GetPermissionByID(id uint) (*models.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *PermissionService) UpdatePermissionByID(id uint, permission *models.Permission) error {
	return s.repo.Update(id, permission)
}

func (s *PermissionService) DeletePermissionByID(id uint) error {
	return s.repo.Delete(id)
}
