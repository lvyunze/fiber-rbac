package service

import (
	"github.com/lvyunze/fiber-rbac/internal/models"
	"github.com/lvyunze/fiber-rbac/internal/repository"
)

type RoleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(role *models.Role) error {
	return s.repo.Create(role)
}

func (s *RoleService) GetRoles() ([]models.Role, error) {
	return s.repo.FindAll()
}

func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.repo.FindByID(id)
}

func (s *RoleService) UpdateRoleByID(id uint, role *models.Role) error {
	return s.repo.Update(id, role)
}

func (s *RoleService) DeleteRoleByID(id uint) error {
	return s.repo.Delete(id)
}
