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
