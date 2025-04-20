package service

import (
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/schema"
)

// PermissionService 权限服务接口
type PermissionService interface {
	Create(req *schema.CreatePermissionRequest) (uint64, error)
	Update(req *schema.UpdatePermissionRequest) error
	Delete(id uint64) error
	GetByID(id uint64) (*schema.PermissionResponse, error)
	List(req *schema.ListPermissionRequest) (*schema.ListPermissionResponse, error)
}

// permissionService 权限服务实现
type permissionService struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

// Create 创建权限
func (s *permissionService) Create(req *schema.CreatePermissionRequest) (uint64, error) {
	// 检查权限名是否已存在
	existingPermission, err := s.permissionRepo.GetByName(req.Name)
	if err != nil {
		return 0, err
	}

	if existingPermission != nil {
		return 0, errors.ErrPermissionExists
	}

	// 检查权限标识是否已存在
	existingPermissionByCode, err := s.permissionRepo.GetByCode(req.Code)
	if err != nil {
		return 0, err
	}

	if existingPermissionByCode != nil {
		return 0, errors.ErrPermissionExists
	}

	// 创建权限
	permission := &model.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return 0, err
	}

	return permission.ID, nil
}

// Update 更新权限
func (s *permissionService) Update(req *schema.UpdatePermissionRequest) error {
	// 检查权限是否存在
	existingPermission, err := s.permissionRepo.GetByID(req.ID)
	if err != nil {
		return err
	}

	if existingPermission == nil {
		return errors.ErrPermissionNotFound
	}

	// 检查权限名是否已被其他权限使用
	if req.Name != existingPermission.Name {
		permission, err := s.permissionRepo.GetByName(req.Name)
		if err != nil {
			return err
		}

		if permission != nil && permission.ID != req.ID {
			return errors.ErrPermissionExists
		}
	}

	// 检查权限标识是否已被其他权限使用
	if req.Code != existingPermission.Code {
		permission, err := s.permissionRepo.GetByCode(req.Code)
		if err != nil {
			return err
		}

		if permission != nil && permission.ID != req.ID {
			return errors.ErrPermissionExists
		}
	}

	// 更新权限信息
	updatedPermission := &model.Permission{
		ID:          req.ID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	return s.permissionRepo.Update(updatedPermission)
}

// Delete 删除权限
func (s *permissionService) Delete(id uint64) error {
	// 检查权限是否存在
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}

	if permission == nil {
		return errors.ErrPermissionNotFound
	}

	// 检查权限是否被角色使用
	if len(permission.Roles) > 0 {
		return errors.ErrPermissionInUse
	}

	// 删除权限
	return s.permissionRepo.Delete(id)
}

// GetByID 根据ID获取权限
func (s *permissionService) GetByID(id uint64) (*schema.PermissionResponse, error) {
	// 获取权限信息
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		return nil, errors.ErrPermissionNotFound
	}

	return s.convertToPermissionResponse(permission), nil
}

// List 获取权限列表，返回完整分页信息
func (s *permissionService) List(req *schema.ListPermissionRequest) (*schema.ListPermissionResponse, error) {
	permissions, total, err := s.permissionRepo.List(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, err
	}

	items := make([]schema.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, *s.convertToPermissionResponse(permission))
	}
	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	return &schema.ListPermissionResponse{
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		Items:      items,
	}, nil
}

// convertToPermissionResponse 将权限模型转换为响应结构
func (s *permissionService) convertToPermissionResponse(permission *model.Permission) *schema.PermissionResponse {
	return &schema.PermissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
	}
}
