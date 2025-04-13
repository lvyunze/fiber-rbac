package repository

import (
	"errors"
	"fmt"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint64) error
	GetByID(id uint64) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetByCode(code string) (*model.Role, error)
	List(page, pageSize int, keyword string) ([]*model.Role, int64, error)
	AddPermissions(roleID uint64, permissionIDs []uint64) error
	RemovePermissions(roleID uint64, permissionIDs []uint64) error
	UpdatePermissions(roleID uint64, permissionIDs []uint64) error
	GetUsersByRoleID(roleID uint64) ([]*model.User, error)
	GetRoleWithPermissions(roleID uint64) (*model.Role, error)
}

// roleRepo 角色仓储实现
type roleRepo struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepo{db: db}
}

// Create 创建角色
func (r *roleRepo) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *roleRepo) Update(role *model.Role) error {
	// 只更新非零值字段
	return r.db.Model(role).Updates(role).Error
}

// Delete 删除角色（软删除）
func (r *roleRepo) Delete(id uint64) error {
	// 使用事务确保数据一致性
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除角色关联的权限
		if err := tx.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 删除用户与角色的关联
		if err := tx.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 软删除角色
		return tx.Model(&model.Role{}).Where("id = ?", id).Update("deleted_at", model.SoftDelete()).Error
	})
}

// GetByID 根据ID获取角色
func (r *roleRepo) GetByID(id uint64) (*model.Role, error) {
	var role model.Role
	result := r.db.Preload("Permissions").Where("id = ? AND deleted_at IS NULL", id).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 角色不存在返回nil，而不是错误
		}
		return nil, result.Error
	}
	return &role, nil
}

// GetByCode 根据编码获取角色
func (r *roleRepo) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	result := r.db.Preload("Permissions").Where("code = ? AND deleted_at IS NULL", code).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

// GetRoleWithPermissions 获取角色及其权限详情
func (r *roleRepo) GetRoleWithPermissions(roleID uint64) (*model.Role, error) {
	var role model.Role
	result := r.db.Preload("Permissions").Where("id = ? AND deleted_at IS NULL", roleID).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

// GetUsersByRoleID 获取拥有指定角色的用户列表
func (r *roleRepo) GetUsersByRoleID(roleID uint64) ([]*model.User, error) {
	// 检查角色是否存在
	exists := r.db.Model(&model.Role{}).Where("id = ? AND deleted_at IS NULL", roleID).First(&model.Role{})
	if exists.Error != nil {
		if errors.Is(exists.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 角色不存在
		}
		return nil, exists.Error
	}

	// 查询拥有该角色的用户
	var users []*model.User
	if err := r.db.Joins("JOIN user_roles ON users.id = user_roles.user_id").Where("user_roles.role_id = ? AND users.deleted_at IS NULL", roleID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetByName 根据名称获取角色
func (r *roleRepo) GetByName(name string) (*model.Role, error) {
	var role model.Role
	result := r.db.Preload("Permissions").Where("name = ? AND deleted_at IS NULL", name).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

// List 获取角色列表
func (r *roleRepo) List(page, pageSize int, keyword string) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := r.db.Model(&model.Role{}).Where("deleted_at IS NULL")

	// 添加关键词搜索
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("Permissions").Offset(offset).Limit(pageSize).Order("id DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AddPermissions 为角色添加权限
func (r *roleRepo) AddPermissions(roleID uint64, permissionIDs []uint64) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查角色是否存在
		var role model.Role
		if err := tx.Where("id = ? AND deleted_at IS NULL", roleID).First(&role).Error; err != nil {
			return err
		}

		// 添加权限关联
		for _, permissionID := range permissionIDs {
			// 检查权限是否存在
			var permission model.Permission
			if err := tx.Where("id = ? AND deleted_at IS NULL", permissionID).First(&permission).Error; err != nil {
				return fmt.Errorf("权限ID %d 不存在: %w", permissionID, err)
			}

			// 检查关联是否已存在
			var count int64
			tx.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Count(&count)
			if count == 0 {
				// 创建关联
				rolePermission := model.RolePermission{
					RoleID:       roleID,
					PermissionID: permissionID,
					CreatedAt:    time.Now().Unix(), // 添加创建时间
				}
				if err := tx.Create(&rolePermission).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// RemovePermissions 移除角色的权限
func (r *roleRepo) RemovePermissions(roleID uint64, permissionIDs []uint64) error {
	return r.db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).Delete(&model.RolePermission{}).Error
}

// UpdatePermissions 更新角色的权限
func (r *roleRepo) UpdatePermissions(roleID uint64, permissionIDs []uint64) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有现有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 添加新权限
		if len(permissionIDs) > 0 {
			return r.AddPermissions(roleID, permissionIDs)
		}

		return nil
	})
}
