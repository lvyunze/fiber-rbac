package service

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/schema"
)

// RoleService 角色服务接口
type RoleService interface {
	Create(req *schema.CreateRoleRequest) (uint64, error)
	Update(req *schema.UpdateRoleRequest) error
	Delete(id uint64) error
	GetByID(id uint64) (*schema.RoleResponse, error)
	List(req *schema.ListRoleRequest) (*schema.ListRoleResponse, error)
	AssignPermission(roleID uint64, permissionIDs []uint64) error
	GetPermissions(roleID uint64) ([]schema.PermissionResponse, error)
}

// roleService 角色服务实现
type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// Create 创建角色
func (s *roleService) Create(req *schema.CreateRoleRequest) (uint64, error) {
	// 检查角色名是否已存在
	existingRole, err := s.roleRepo.GetByName(req.Name)
	if err != nil {
		return 0, err
	}

	if existingRole != nil {
		return 0, errors.ErrRoleExists
	}

	// 检查角色编码是否已存在
	existingRoleByCode, err := s.roleRepo.GetByCode(req.Code)
	if err != nil {
		return 0, err
	}

	if existingRoleByCode != nil {
		return 0, errors.ErrRoleExists
	}

	// 创建角色
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return 0, err
	}

	// 如果有指定权限，添加权限关联
	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepo.AddPermissions(role.ID, req.PermissionIDs); err != nil {
			slog.Error("添加角色权限失败", "error", err)
			// 不返回错误，继续执行
		}
	}

	return role.ID, nil
}

// Update 更新角色
func (s *roleService) Update(req *schema.UpdateRoleRequest) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(req.ID)
	if err != nil {
		return err
	}

	if existingRole == nil {
		return errors.ErrRoleNotFound
	}

	// 检查角色名是否已被其他角色使用
	if req.Name != existingRole.Name {
		role, err := s.roleRepo.GetByName(req.Name)
		if err != nil {
			return err
		}

		if role != nil && role.ID != req.ID {
			return errors.ErrRoleExists
		}
	}

	// 检查角色标识是否已被其他角色使用
	if req.Code != existingRole.Code {
		role, err := s.roleRepo.GetByCode(req.Code)
		if err != nil {
			return err
		}

		if role != nil && role.ID != req.ID {
			return errors.ErrRoleExists
		}
	}

	// 更新角色信息
	updatedRole := &model.Role{
		ID:          req.ID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := s.roleRepo.Update(updatedRole); err != nil {
		return err
	}

	// 如果提供了权限ID，更新角色权限
	if req.PermissionIDs != nil {
		if err := s.roleRepo.UpdatePermissions(req.ID, req.PermissionIDs); err != nil {
			slog.Error("更新角色权限失败", "error", err)
			// 不返回错误，继续执行
		}
	}

	return nil
}

// Delete 删除角色
func (s *roleService) Delete(id uint64) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}

	if role == nil {
		return errors.ErrRoleNotFound
	}

	// 检查角色是否被用户使用
	users, err := s.roleRepo.GetUsersByRoleID(id)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.ErrRoleInUse
	}

	// 删除角色
	return s.roleRepo.Delete(id)
}

// GetByID 根据ID获取角色
func (s *roleService) GetByID(id uint64) (*schema.RoleResponse, error) {
	// 获取角色信息
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.ErrRoleNotFound
	}

	return s.convertToRoleResponse(role), nil
}

// List 获取角色列表，返回完整分页信息
func (s *roleService) List(req *schema.ListRoleRequest) (*schema.ListRoleResponse, error) {
	orderBy := req.OrderBy
	desc := req.Desc
	roles, total, err := s.roleRepo.List(req.Page, req.PageSize, req.Keyword, orderBy, desc)
	if err != nil {
		return nil, err
	}

	items := make([]schema.RoleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, *s.convertToRoleResponse(role))
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	return &schema.ListRoleResponse{
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		TotalPages:  totalPages,
		Items:       items,
	}, nil
}

// AssignPermission 分配权限给角色
func (s *roleService) AssignPermission(roleID uint64, permissionIDs []uint64) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	if role == nil {
		return errors.ErrRoleNotFound
	}

	// 检查所有权限是否存在
	for _, permID := range permissionIDs {
		perm, err := s.permissionRepo.GetByID(permID)
		if err != nil {
			return err
		}
		if perm == nil {
			return errors.ErrPermissionNotFound
		}
	}

	// 更新角色权限
	return s.roleRepo.UpdatePermissions(roleID, permissionIDs)
}

// GetPermissions 获取角色的权限列表
func (s *roleService) GetPermissions(roleID uint64) ([]schema.PermissionResponse, error) {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.ErrRoleNotFound
	}

	// 获取并加载权限
	role, err = s.roleRepo.GetRoleWithPermissions(roleID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	permissions := make([]schema.PermissionResponse, 0, len(role.Permissions))
	for _, perm := range role.Permissions {
		permissions = append(permissions, schema.PermissionResponse{
			ID:          perm.ID,
			Code:        perm.Code,
			Name:        perm.Name,
			Description: perm.Description,
			CreatedAt:   perm.CreatedAt,
		})
	}

	return permissions, nil
}

// convertToRoleResponse 将角色模型转换为响应结构
func (s *roleService) convertToRoleResponse(role *model.Role) *schema.RoleResponse {
	response := &schema.RoleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		Permissions: make([]schema.PermissionSimple, 0, len(role.Permissions)),
	}

	// 添加权限信息
	for _, perm := range role.Permissions {
		response.Permissions = append(response.Permissions, schema.PermissionSimple{
			ID:   perm.ID,
			Code: perm.Code,
			Name: perm.Name,
		})
	}

	return response
}
