package repository

import (
	"errors"
	"fmt"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"strings"

	"gorm.io/gorm"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	Create(permission *model.Permission) error
	Update(permission *model.Permission) error
	Delete(id uint64) error
	GetByID(id uint64) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	GetByCode(code string) (*model.Permission, error)
	List(page, pageSize int, keyword string) ([]*model.Permission, int64, error)
	GetByNames(names []string) ([]*model.Permission, error)
}

// permissionRepo 权限仓储实现
type permissionRepo struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepo{db: db}
}

// Create 创建权限
func (r *permissionRepo) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// Update 更新权限
func (r *permissionRepo) Update(permission *model.Permission) error {
	// 只更新非零值字段
	return r.db.Model(permission).Updates(permission).Error
}

// Delete 删除权限（软删除）
func (r *permissionRepo) Delete(id uint64) error {
	// 使用事务确保数据一致性
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除角色关联的权限
		if err := tx.Where("permission_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 软删除权限
		return tx.Model(&model.Permission{}).Where("id = ?", id).Update("deleted_at", model.SoftDelete()).Error
	})
}

// GetByID 根据ID获取权限
func (r *permissionRepo) GetByID(id uint64) (*model.Permission, error) {
	var permission model.Permission
	result := r.db.Preload("Roles").Where("id = ? AND deleted_at IS NULL", id).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 权限不存在返回nil，而不是错误
		}
		return nil, result.Error
	}
	return &permission, nil
}

// GetByName 根据名称获取权限
func (r *permissionRepo) GetByName(name string) (*model.Permission, error) {
	var permission model.Permission
	result := r.db.Preload("Roles").Where("name = ? AND deleted_at IS NULL", name).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &permission, nil
}

// GetByCode 根据编码获取权限
func (r *permissionRepo) GetByCode(code string) (*model.Permission, error) {
	var permission model.Permission
	result := r.db.Preload("Roles").Where("code = ? AND deleted_at IS NULL", code).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &permission, nil
}

// List 获取权限列表
func (r *permissionRepo) List(page, pageSize int, keyword string) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := r.db.Model(&model.Permission{}).Where("deleted_at IS NULL")

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
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetByNames 根据权限名称列表获取权限
func (r *permissionRepo) GetByNames(names []string) ([]*model.Permission, error) {
	var permissions []*model.Permission
	if len(names) == 0 {
		return permissions, nil
	}

	if err := r.db.Where("name IN ? AND deleted_at IS NULL", names).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}
